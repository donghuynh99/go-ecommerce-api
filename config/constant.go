package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	PostgresqlConfig PostgresqlConfig
	AppConfig        AppConfig
}

type PostgresqlConfig struct {
	Database string `envconfig:"POSTGRESQL_DB" required:"false" defaul:"ecommerce_api"`
	Host     string `envconfig:"POSTGRESQL_HOST" required:"false" default:"postgresql"`
	Port     int    `envconfig:"POSTGRESQL_PORT" required:"false" default:"5432"`
	User     string `envconfig:"POSTGRESQL_USER" required:"false" default:"root"`
	Password string `envconfig:"POSTGRESQL_PASSWORD" required:"false" default:"pass"`
}

type AppConfig struct {
	AppName string `envconfig:"APP_NAME" required:"false" default:"Ecommerce API"`
	AppPort int    `envconfig:"APP_PORT" required:"false" default:"8080"`
	AppURL  string `envconfig:"APP_PORT" required:"false" default:"http://localhost:8080"`
}

var c *Config

func Load() *Config {
	var cnf Config
	if c != nil {
		return c
	}
	err := envconfig.Process("", &cnf)
	if err != nil {
		panic(err)
	}
	c = &cnf
	return c
}

func GetConfig() *Config {
	return c
}
