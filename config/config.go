package config

import (
	"fmt"
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
		return nil, err
	}

	config:= &Config{
		TgToken: os.Getenv("TgToken"),
	}

	if config.TgToken == "" {
		return nil, fmt.Errorf("empty env token")
	}

	return config, nil
}
