package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	WebappBaseUrl string `mapstructure:"WEBAPP_BASE_URL"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	EnvType       string `mapstructure:"ENV_TYPE"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	GoogleAPI     string `mapstructure:"GOOGLE_API"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + ".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
	}

	return
}
