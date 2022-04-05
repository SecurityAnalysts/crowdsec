package v1

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/crowdsecurity/crowdsec/pkg/database"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/machine"
	"github.com/crowdsecurity/crowdsec/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var identityKey = "id"

type JWT struct {
	Middleware *jwt.GinJWTMiddleware
	DbClient   *database.Client
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if value, ok := data.(*models.WatcherAuthRequest); ok {
		return jwt.MapClaims{
			identityKey: &value.MachineID,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	machineId := claims[identityKey].(string)
	return &models.WatcherAuthRequest{
		MachineID: &machineId,
	}
}

var AllowedOU = "bouncer"

func (just *JWT) ValidateCert(c *gin.Context) (bool, string) {
	if c.Request.TLS != nil {
		if len(c.Request.TLS.VerifiedChains) > 0 {
			validOU := false
			clientCert := c.Request.TLS.VerifiedChains[0][0]
			for _, ou := range clientCert.Subject.OrganizationalUnit {
				if AllowedOU == ou {
					validOU = true
				}
			}
			if !validOU {
				log.Errorf("APIKey: client certificate is not from an agent")
				return false, ""
			}
			return true, fmt.Sprintf("%s-%s", clientCert.Subject.CommonName, c.ClientIP())
		} else {
			log.Error("Found no verified certs in request")
			return false, ""
		}
	}
	return false, ""
}

func (j *JWT) Authenticator(c *gin.Context) (interface{}, error) {
	var loginInput models.WatcherAuthRequest
	var scenarios string
	var err error
	var scenariosInput []string
	var clientMachine *ent.Machine
	var machineID string

	if c.Request.TLS != nil && len(c.Request.TLS.PeerCertificates) > 0 {
		validCert, machineID := j.ValidateCert(c)
		if !validCert {
			c.JSON(http.StatusForbidden, gin.H{"message": "access forbidden"})
			c.Abort()
			return nil, fmt.Errorf("failed cert authentication")
		}
		log.Printf("certificate auth ok for agent %s", machineID)

		clientMachine, err = j.DbClient.Ent.Machine.Query().
			Where(machine.MachineId(machineID)).
			First(j.DbClient.CTX)
		if ent.IsNotFound(err) {
			//Machine was not found, let's create it
			log.Printf("machine %s not found, create it", machineID)
			password := strfmt.Password("")
			_, err = j.DbClient.CreateMachine(&machineID, &password, "", true, true)
			if err != nil {
				return "", errors.Wrapf(err, "while creating machine entry for %s", machineID)
			}
			clientMachine, err = j.DbClient.Ent.Machine.Query().
				Where(machine.MachineId(machineID)).
				First(j.DbClient.CTX)
			if err != nil {
				return "", errors.Wrapf(err, "while selecting machine entry for %s after creation", machineID)
			}
		} else {
			machineID = clientMachine.MachineId
			//we should still get the updated list of scenarios from the machine
		}

	} else {
		//normal auth

		if err := c.ShouldBindJSON(&loginInput); err != nil {
			return "", errors.New(fmt.Sprintf("missing : %v", err.Error()))
		}
		if err := loginInput.Validate(strfmt.Default); err != nil {
			return "", errors.New("input format error")
		}
		machineID = *loginInput.MachineID
		password := *loginInput.Password
		scenariosInput = loginInput.Scenarios

		clientMachine, err = j.DbClient.Ent.Machine.Query().
			Where(machine.MachineId(machineID)).
			First(j.DbClient.CTX)
		if err != nil {
			log.Printf("Error machine login for %s : %+v ", machineID, err)
			return nil, err
		}

		if clientMachine == nil {
			log.Errorf("Nothing for '%s'", machineID)
			return nil, jwt.ErrFailedAuthentication
		}

		if !clientMachine.IsValidated {
			return nil, fmt.Errorf("machine %s not validated", machineID)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(clientMachine.Password), []byte(password)); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		//end of normal auth
	}

	if len(scenariosInput) > 0 {
		for _, scenario := range scenariosInput {
			if scenarios == "" {
				scenarios = scenario
			} else {
				scenarios += "," + scenario
			}
		}
		err = j.DbClient.UpdateMachineScenarios(scenarios, clientMachine.ID)
		if err != nil {
			log.Errorf("Failed to update scenarios list for '%s': %s\n", machineID, err)
			return nil, jwt.ErrFailedAuthentication
		}
	}

	if clientMachine.IpAddress == "" {
		err = j.DbClient.UpdateMachineIP(c.ClientIP(), clientMachine.ID)
		if err != nil {
			log.Errorf("Failed to update ip address for '%s': %s\n", machineID, err)
			return nil, jwt.ErrFailedAuthentication
		}
	}

	if clientMachine.IpAddress != c.ClientIP() && clientMachine.IpAddress != "" {
		log.Warningf("new IP address detected for machine '%s': %s (old: %s)", clientMachine.MachineId, c.ClientIP(), clientMachine.IpAddress)
		err = j.DbClient.UpdateMachineIP(c.ClientIP(), clientMachine.ID)
		if err != nil {
			log.Errorf("Failed to update ip address for '%s': %s\n", clientMachine.MachineId, err)
			return nil, jwt.ErrFailedAuthentication
		}
	}

	useragent := strings.Split(c.Request.UserAgent(), "/")
	if len(useragent) != 2 {
		log.Warningf("bad user agent '%s' from '%s'", c.Request.UserAgent(), c.ClientIP())
		return nil, jwt.ErrFailedAuthentication
	}

	if err := j.DbClient.UpdateMachineVersion(useragent[1], clientMachine.ID); err != nil {
		log.Errorf("unable to update machine '%s' version '%s': %s", clientMachine.MachineId, useragent[1], err)
		log.Errorf("bad user agent from : %s", c.ClientIP())
		return nil, jwt.ErrFailedAuthentication
	}

	return &models.WatcherAuthRequest{
		MachineID: &machineID,
	}, nil

}

func Authorizator(data interface{}, c *gin.Context) bool {
	return true
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func NewJWT(dbClient *database.Client) (*JWT, error) {
	// Get secret from environment variable "SECRET"
	var (
		secret []byte
	)

	//Please be aware that brute force HS256 is possible.
	//PLEASE choose a STRONG secret
	secret_string := os.Getenv("CS_LAPI_SECRET")
	if secret_string == "" {
		secret = make([]byte, 64)
		if n, err := rand.Read(secret); err != nil {
			log.Fatalf("unable to generate a new random seed for JWT generation")
		} else {
			if n != 64 {
				log.Fatalf("not enough entropy at random seed generation for JWT generation")
			}
		}
	} else {
		secret = []byte(secret_string)
		if len(secret) < 64 {
			log.Fatalf("secret not strong enough")
		}
	}

	jwtMiddleware := &JWT{
		DbClient: dbClient,
	}

	ret, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Crowdsec API local",
		Key:             secret,
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
		Authenticator:   jwtMiddleware.Authenticator,
		Authorizator:    Authorizator,
		Unauthorized:    Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	errInit := ret.MiddlewareInit()
	if errInit != nil {
		return &JWT{}, fmt.Errorf("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	if err != nil {
		return &JWT{}, err
	}

	return &JWT{Middleware: ret}, nil
}
