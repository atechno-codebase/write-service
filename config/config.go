package config

import (
	"github.com/AbsaOSS/env-binder/env"
	"github.com/tidwall/gjson"
)

var configuration gjson.Result

type Config struct {
	LogPath     string `env:"LOG_PATH"`
	SecretToken string `env:"SECRET_TOKEN"`
	Port        string `env:"PORT,default=8000"`
	MongoURL    string `env:"MONGO_URL"`
	Database    string `env:"DATABASE"`
}

var Configuration Config

func Init() error {
	return env.Bind(&Configuration)
}
