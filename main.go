package main

import (
	"github.com/library/app"
	"github.com/library/mail"
	"github.com/spf13/viper"
	"log"
	"os"
)
// @title           Library API
// @version         2.0
// @description     This is a server API.
// @host            localhost:8000
// @securityDefinitions.apikey mySigningKey
// @in header
// @name Authorization

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// Log current environment.
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	log.Println("ENV: " + current_env)

	a := app.App{}

	a.Initialize()
	if os.Getenv("PORT") == "" {
		// Get port from config if no env variable.
		a.Run(":" + viper.GetString("PORT"))
	} else {
		// Get port from env.
		a.Run(":" + os.Getenv("PORT"))
	}

	email := mail.NewEmail("xxx@126.com", "golang mail", "please, return books to the library")
	err = mail.SendEmail(email)
	log.Print(err)

}
