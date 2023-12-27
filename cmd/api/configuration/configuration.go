package configuration

import (
	"github.com/ecumenos/fxecumenos/fxlogger"
	"github.com/ecumenos/fxecumenos/fxpostgres"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/orbis-socius-register/cmd/api/httpserver"
	"github.com/ecumenos/orbis-socius-register/cmd/api/smtp"
)

type Config struct {
	APILocal           bool `default:"true"`
	APILogger          *fxlogger.Config
	APIHTTP            *httpserver.Config
	APIDataStore       *fxpostgres.Config
	APIResponseFactory *fxrf.Config
	APISMTPSender      *smtp.SenderConfig
}
