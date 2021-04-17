package main

import (
	"github.com/spf13/viper"
)

var (
	ConfigName      string    = "local"
	ConfigExtension string    = "yaml"
	ConfigPath      [2]string = [2]string{".", "./config"}
	EnvVarPrefix    string    = "ENV"
)

type Config struct {
	DB_HOST     string
	DB_TYPE     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string

	SERVE_ADDR string
}

// Uses viper lib to read config file and env variables
// Env variables have precedence
// To read env vars, a config file must be read (hence why there exists ./skeleton.yaml)
func InitializeConfig() (*Config, error) {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigExtension)

	for _, path := range ConfigPath {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix(EnvVarPrefix)
	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
