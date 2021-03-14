package config

import "github.com/spf13/viper"

func Setup() {
	viper.SetDefault("DbHost", "localhost")
	viper.SetDefault("DbName", "messaging-service")
	viper.SetDefault("DbUser", "messaging-service")
	viper.SetDefault("DbPass", "123456")
	viper.SetDefault("DbPort", "5432")
}
