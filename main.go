package main

import (
	"github.com/library/app"
	"github.com/library/callAt"
	"github.com/library/db"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(time.Hour)
	task := make(chan []string)

	go func() {
		for {
			select {
			case <-ticker.C:
				listEmail, err := callAt.CheckReturnDate(db.DB{})
				if err != nil {
					log.Printf("Can not check return data for issue acts (%s):%s", time.Now(), err)
				}
				if len(listEmail) > 0 {
					task <- listEmail
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-task:
				listEmail := <-task
				for _, v := range listEmail{
					callAt.Email([]string{v})
				}
				log.Println("Email Sent Successfully!")
			}
		}

	}()

	<-quit
	ticker.Stop()
}
