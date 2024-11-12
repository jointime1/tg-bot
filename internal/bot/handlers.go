package bot

import (
	"fmt"

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
	if msg.NewChatMembers != nil {
        deleteMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
        b.api.DeleteMessage(deleteMsg)

        for _, member := range *msg.NewChatMembers {
            welcomeMsg := tgbotapi.NewMessage(msg.Chat.ID, "Добро пожаловать, в IT-ХАЗЯЕВА!"+member.FirstName+"!")
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




}

// func (b *Bot) handleUnknownCommand(msg *tgbotapi.Message) {
// 	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unknown command"))
// }


// func (b *Bot) handleTextMessage(msg *tgbotapi.Message) {
// 	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, msg.Text))
// }

// func (b *Bot) handleStartCommand(msg *tgbotapi.Message) {
// 	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Hello, "+msg.Chat.FirstName))
// }