package database

import (
	"fmt"
	"github.com/library/src/config"
	"github.com/library/src/models"
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

	db.AutoMigrate(&models.Book{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Lending{})
	db.AutoMigrate(&models.Acceptance{})

	DB = db
}
