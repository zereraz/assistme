package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zereraz/assistme/config"
	"github.com/zereraz/assistme/db"
	l "github.com/zereraz/assistme/log"
	"github.com/zereraz/assistme/message"
)

func init() {
	_, err := db.SetupDb()
	if err != nil {
		panic(fmt.Sprintf("Could not setup db: %v", err))
	}
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Cleaning up")
		db.Cleanup()
		os.Exit(0)
	}()
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.API_TOKEN)
	if err != nil {
		l.Log.Fatalf("Could not setup telegram bot %v", err)
	}
	setupCloseHandler()
	message.ListenToCommands(bot)
}
