package main

import (
	"bot/config"
	"bot/internal/auth"
	"bot/internal/bot"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joeyak/go-twitch-eventsub/v3"
)

type TwitchStreamEvent struct {
    UserName   string `json:"user_name"`  // Имя стримера
    Title      string `json:"title"`      // Название стрима
    GameName   string `json:"game_name"`  // Название игры/категории
    ViewerCount int   `json:"viewer_count"` // Количество зрителей
    StartedAt  string `json:"started_at"` // Время начала стрима
}



func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Ошибка конфигурации: %v", err)
	}

	bot, err := bot.NewBot(config.TgToken)
	if err != nil {
		log.Printf("Ошибка инициализации бота: %v", err)
	}

	userTokenChan := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		token, err := auth.GetTwitchUserToken(config, code, config.RedirectURI)
		if err != nil {
			log.Printf("Ошибка получения user access token: %v", err)
			http.Error(w, "Failed to get user token", http.StatusInternalServerError)
			return
		}

		userTokenChan <- token

		fmt.Fprintf(w, "Авторизация успешна. User Access Token: %s", token)
	})

	go func() {
		log.Println("Сервер запущен на http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	go func() {
			log.Println("Запуск бота...")
			
			bot.Start()
			
	}()

	userToken := <-userTokenChan

	userID, err := auth.GetTwitchUserId(config, userToken)
	if err != nil {
		log.Printf("Ошибка получения user id: %v", err)
	}
	fmt.Println(userID)

	

	// Вебсокеты
	client := twitch.NewClient()

	client.OnError(func(err error) {
		fmt.Printf("ERROR: %v\n", err)
	})
	client.OnWelcome(func(message twitch.WelcomeMessage) {
		fmt.Printf("WELCOME: %v\n", message)

		events := []twitch.EventSubscription{
			twitch.SubStreamOnline,
			twitch.SubStreamOffline,
		}

		for _, event := range events {
			fmt.Printf("subscribing to %s\n", event)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    config.ClientId,
				AccessToken: userToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": userID,
				},
			})
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
		}
	})
	client.OnNotification(func(message twitch.NotificationMessage) {
		if(message.Payload.Subscription.Type == "stream.online") {
			//ChatId беседы, куда срать уведомлениями
			fmt.Println(message.Payload.Subscription)
			var event TwitchStreamEvent
			err := json.Unmarshal([]byte(*message.Payload.Event), &event)
				if err != nil {
					fmt.Printf("Ошибка при разборе данных стрима: %v\n", err)
					return
				}

			broadcasterName := event.UserName        // Имя стримера
			streamTitle := event.Title               // Название стрима
			streamCategory := event.GameName        // Категория игры, в которую играет стример
			streamViewerCount := event.ViewerCount  // Количество зрителей
			streamStartTime := event.StartedAt 

			fmt.Printf("%s: %s, %s, %d, %s\n", broadcasterName, streamTitle, streamCategory, streamViewerCount, streamStartTime)
			
			chatID:=-4551123737
			bot.SendStreamOnlineMessage(int64(chatID))
		}
		fmt.Printf("NOTIFICATION: %s: %#v\n", message.Payload.Subscription.Type, message.Payload.Event)
	})
	client.OnKeepAlive(func(message twitch.KeepAliveMessage) {
		fmt.Printf("KEEPALIVE: %v\n", message)
	})
	client.OnRevoke(func(message twitch.RevokeMessage) {
		fmt.Printf("REVOKE: %v\n", message)
	})
	client.OnRawEvent(func(event string, metadata twitch.MessageMetadata, subscription twitch.PayloadSubscription) {
		fmt.Printf("EVENT[%s]: %s: %s\n", subscription.Type, metadata, event)
	})

	err = client.Connect()
	if err != nil {
		fmt.Printf("Could not connect client: %v\n", err)
	}


}
