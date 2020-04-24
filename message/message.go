package message

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	l "github.com/raunaqrox/assistme/log"
)

func ListenToCommands(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "") //update.Message.Text)
		fmt.Println(update.Message.Chat.ID)
		if update.Message.IsCommand() {
			l.Log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("Got command %s", update.Message)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				l.Log.Fatalf("Error while sending message: %v", err)
			}
		} else {
			continue
		}
	}
}
