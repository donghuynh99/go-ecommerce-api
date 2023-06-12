package database

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/donghuynh99/ecommerce_api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Connect() {
	conf := config.GetConfig().PostgresqlConfig
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.User, conf.Password, conf.DB, strconv.Itoa(conf.Port))

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	d, _ := db.DB()
	d.SetMaxIdleConns(100)
	d.SetConnMaxIdleTime(500 * time.Microsecond)

	log.Println("Connect Postgres OK")

	Database = db
}
