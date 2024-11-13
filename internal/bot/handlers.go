package bot

import (
	"bot/config"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) Start()  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message != nil { 
			b.handleMessage(update.Message)
		}
	}

}



func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	fmt.Println("zxc")
	if msg.NewChatMembers != nil {
        deleteMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
        b.api.DeleteMessage(deleteMsg)

        for _, member := range *msg.NewChatMembers {
			welcomeMsg := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("<b>%s</b> присоединился к чату", member.FirstName))
			welcomeMsg.ParseMode = "HTML"
			
            b.api.Send(welcomeMsg)
        }
        return
    }

	if msg.LeftChatMember != nil {	
		deleteMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
		
		if _, err:= b.api.DeleteMessage(deleteMsg);err != nil {
			fmt.Println(msg.Chat.ID, msg.MessageID)
			return
		}

		return
	}

	if msg.Text == "/auth" {
		fmt.Println("auth")
		url:= authURL()
		authMsg := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("<a href=\"%s\">Авторизация</a>", url))
		authMsg.ParseMode = "HTML"
		
		b.api.Send(authMsg)
	}
}


func authURL() string {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Ошибка конфигурации: %v", err)
	}
	fmt.Println(config.ClientId)
    return fmt.Sprintf(
        "https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=user:read:email",
        config.ClientId,
        config.RedirectURI,
    )
}

func (b *Bot) SendStreamOnlineMessage(chatID int64) {
	message := tgbotapi.NewMessage(chatID, "Стрим начался!")
	b.api.Send(message)
}
