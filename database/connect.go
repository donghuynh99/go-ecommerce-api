package database

import (
	"fmt"
	"log"
	"time"

	"github.com/donghuynh99/ecommerce_api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Connect() {
	conf := config.GetConfig().PostgresqlConfig
	dns := fmt.Sprintf("postgres://%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("Connect Postgres Fail!")
	}

	d, _ := db.DB()
	d.SetMaxIdleConns(100)
	d.SetConnMaxIdleTime(500 * time.Microsecond)

	log.Println("Connect Postgres OK")

	Database = db
}
