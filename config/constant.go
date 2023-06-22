package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/shopspring/decimal"
)

type Config struct {
	PostgresqlConfig  PostgresqlConfig
	AppConfig         AppConfig
	RoleConfig        RoleConfig
	PaginationConfig  PaginationConfig
	GeneralConfig     GeneralConfig
	StatusOrderConfig StatusOrderConfig
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

type PaginationConfig struct {
	Limit string `default:"10"`
	Page  string `default:"1"`
}

type GeneralConfig struct {
	ImageLimit                   int    `default:"5"`
	DestinationStoreProductImage string `default:"assets/products/images/"`
	DestinationStoreAvatarUser   string `default:"assets/avatar/"`
	KeyToken                     string `default:"private"`
}

type StatusOrderConfig struct {
	Pending   string `default:"pending"`
	Approved  string `default:"approved"`
	Completed string `default:"completed"`
	Canceled  string `default:"canceled"`
}

type AppConfig struct {
	AppName          string `envconfig:"APP_NAME" required:"false" default:"Ecommerce API"`
	AppPort          int    `envconfig:"APP_PORT" required:"false" default:"8080"`
	AppURL           string `envconfig:"APP_URL" required:"false" default:"http://localhost:8080"`
	DefaultImageURL  string `envconfig:"DEFAULT_IMAGE_URL" required:"false" default:"assets/default_image.jpeg"`
	DefaultAvatarURL string `envconfig:"DEFAULT_AVATAR_URL" required:"false" default:"assets/default_avatar.jpeg"`
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
	ProductName string          `json:"product_name"`
	Thumbnail   string          `json:"product_thumbnail"`
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
}

type ProductJsonStruct struct {
	ID           uint            `json:"id"`
	Name         string          `json:"name"`
	Price        decimal.Decimal `json:"price"`
	Description  string          `json:"description"`
	ThumbnailURL ImageStruct     `json:"thumbnail_url"`
}

type OrderJsonStruct struct {
	ID               uint                  `json:"id"`
	Address          AddressJsonStruct     `json:"address"`
	OrderItems       []OrderItemJsonStruct `json:"order_items"`
	Note             string                `json:"note"`
	ApprovedBy       *uint                 `gorm:"size:255;index"`
	ApprovedAt       *time.Time            `json:"approved_at"`
	CompletedAt      *time.Time            `json:"completed_at"`
	CancelledBy      *uint                 `json:"cancelled_by"`
	CancelledAt      *time.Time            `json:"cancelled_at"`
	CancellationNote *string               `json:"cancellation_note"`
}

type AddressJsonStruct struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	PostCode string `json:"post_code"`
}

type OrderItemJsonStruct struct {
	Name         string          `json:"name"`
	ThumbnailURL ImageStruct     `json:"thumbnail_url"`
	Price        decimal.Decimal `json:"price"`
	Quantity     int             `json:"quantity"`
}

type ImageStruct struct {
	Path string `json:"path"`
	Alt  string `json:"alt"`
}

type AddressListStruct struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsPrimary bool   `json:"is_primary"`
	PostCode  string `json:"post_code"`
}

type ProfileStruct struct {
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Avatar    ImageStruct         `json:"avatar"`
	Email     string              `json:"email"`
	Addresses []AddressListStruct `json:"addresses"`
}
