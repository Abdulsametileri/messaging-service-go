package config

import (
	"github.com/spf13/viper"
	"log"
)

func Setup() {
	if IsDebug {
		viper.Set("DB_HOST", "localhost")
		viper.Set("DB_NAME", "messaging-service")
		viper.Set("DB_USER", "messaging-service")
		viper.Set("DB_PASSWORD", "123456")
		viper.Set("DB_PORT", "5432")

		viper.Set("REDIS_URL", "localhost:6379")
		viper.Set("REDIS_PASSWORD", "samet123")
	} else {
		viper.SetConfigFile(".env")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error while reading config file %s", err)
		}
	}
}
