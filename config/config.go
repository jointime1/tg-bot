package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
type Config struct {
	TgToken string
	UserId string
	ClientSecret string
	ClientId string
	RedirectURI string	
}



// GetConfig returns the config
func GetConfig() (*Config, error) {

	fmt.Println("start")


	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config:= &Config{
		TgToken: os.Getenv("TgToken"),
		UserId: os.Getenv("UserId"),
		ClientSecret: os.Getenv("ClientSecret"),
		ClientId: os.Getenv("ClientId"),
		RedirectURI: "https://b88c-46-242-9-50.ngrok-free.app/callback",
	}

	if config.TgToken == "" {
		return nil, fmt.Errorf("empty env token")
	}

	return config, nil
}
