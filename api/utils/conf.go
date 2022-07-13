package utils

import (
	"fmt"

	env "github.com/Netflix/go-env"
)

type Conf struct {
	// required
	AllowedOrigins         string `env:"ALLOWED_ORIGINS,required=true"`
	DBName                 string `env:"MONGO_INITDB_DATABASE,required=true"`
	DBHostname             string `env:"DB_HOSTNAME,required=true"`
	DBUsername             string `env:"MONGO_INITDB_ROOT_USERNAME,required=true"`
	DBPassword             string `env:"MONGO_INITDB_ROOT_PASSWORD,required=true"`
	GoogleClientID         string `env:"GOOGLE_CLIENT_ID,required=true"`
	SchemaVersion          int    `env:"SCHEMA_VERSION,required=true"`
	SessionLifetimeSeconds int    `env:"SESSION_LIFETIME_SECONDS,required=true"`
	// optional
	ApiLogFile   string `env:"API_LOG_FILE,default=/var/log/linkshare.log"`
	DBPort       string `env:"DB_PORT,default=27017"`
	IsProduction bool   `env:"IS_PRODUCTION,default=true"`
	DebugMode    bool   `env:"DEBUG_MODE,default=false"`
	FrontendURL  string `env:"FRONTEND_URL,default=localhost:3000"`
	// derived
	ConnectionURL string
	scheme        string
}

var gConf *Conf

func GetConf() *Conf {
	if gConf != nil {
		return gConf
	}
	gConf = new(Conf)
	_, err := env.UnmarshalFromEnviron(gConf)
	if err != nil {
		panic(err)
	}
	gConf.SetConnectionURL()
	if gConf.IsProduction {
		gConf.DebugMode = false
		gConf.scheme = "https://"
	} else {
		gConf.AllowedOrigins = "*"
		gConf.scheme = "http://"
	}
	gConf.FrontendURL = gConf.scheme + gConf.FrontendURL

	return gConf
}

// set the connection URL based on db parameters
func (conf *Conf) SetConnectionURL() {
	conf.ConnectionURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
		conf.DBUsername, conf.DBPassword, conf.DBHostname, conf.DBPort)
}
