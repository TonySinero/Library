package models

import (
	"fmt"
	"github.com/library/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB instance

var DB *gorm.DB

// ConnectDataBase return connection

func ConnectDataBase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		config.PostgresHostName,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDataBase,
		config.PostgresPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&User{})

	DB = db
}
