package conf

import (
	"fmt"

	env "github.com/Netflix/go-env"
)

type Conf struct {
	DBName        string `env:"MONGO_INITDB_DATABASE"`
	Hostname      string `env:"HOSTNAME"`
	DBUsername    string `env:"MONGO_INITDB_ROOT_USERNAME,required=true"`
	DBPassword    string `env:"MONGO_INITDB_ROOT_PASSWORD,required=true"`
	SchemaVersion int    `env:"SCHEMA_VERSION,required=true"`
	ConnectionURL string
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
	gConf.ConnectionURL = fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=admin",
		gConf.DBUsername, gConf.DBPassword, gConf.Hostname)
	return gConf
}
