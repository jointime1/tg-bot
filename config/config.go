package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
type Config struct {
	TgToken string	
}



// GetConfig returns the config
func GetConfig() (*Config, error) {

	fmt.Println("start")


	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config:= &Config{
		TgToken: os.Getenv("TgToken"),
	}

	if config.TgToken == "" {
		return nil, fmt.Errorf("empty env token")
	}

	return config, nil
}
