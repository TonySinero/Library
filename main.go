package main

import (
	"log"
	"os"

	"github.com/library/app"
	"github.com/spf13/viper"
)
// @title           Library API
// @version         2.0
// @description     This is a sample server API.
// @host            localhost:8080
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
}
