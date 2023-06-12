package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	PostgresqlConfig PostgresqlConfig
	AppConfig        AppConfig
	RoleConfig       RoleConfig
}

type PostgresqlConfig struct {
	DB       string `envconfig:"POSTGRESQL_DB" required:"false" default:"ecommerce_api"`
	Host     string `envconfig:"POSTGRESQL_HOST" required:"false" default:"localhost"`
	Port     int    `envconfig:"POSTGRESQL_PORT" required:"false" default:"5432"`
	User     string `envconfig:"POSTGRESQL_USER" required:"false" default:"root"`
	Password string `envconfig:"POSTGRESQL_PASSWORD" required:"false" default:"pass"`
}

type RoleConfig struct {
	Admin string `default:"admin"`
	User  string `default:"user"`
}

type AppConfig struct {
	AppName         string `envconfig:"APP_NAME" required:"false" default:"Ecommerce API"`
	AppPort         int    `envconfig:"APP_PORT" required:"false" default:"8080"`
	AppURL          string `envconfig:"APP_URL" required:"false" default:"http://localhost:8080"`
	DefaultImageURL string `envconfig:"DEFAULT_IMAGE_URL" required:"false" default:"assets/default_image.jpeg"`
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

type PaginationStruct struct {
	Data       interface{} `json:"data"`
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
}

type CartInformation struct {
	ProductName string `json:"product_name"`
	Thumbnail   string `json:"product_thumbnail"`
	Quantity    int    `json:"quantity"`
}
