package main

import (
	"bot/config"
	"bot/internal/bot"
	"log"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}


	bot, err := bot.NewBot(config.TgToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	bot.Start()
	
}