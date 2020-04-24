package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Could not load dotenv file: %+v\n", err))
	}

	TelegramBotChannel, err = strconv.Atoi(os.Getenv("BOT_CHANNEL_ID"))
	if err != nil {
		panic("Could not parse Bot channel id")
	}

}

var (
	ENV                = os.Getenv("ENV")
	API_TOKEN          = os.Getenv("TELEGRAM_APITOKEN")
	APP_NAME           = "assistment"
	LogPath            = fmt.Sprintf("/var/log/%s/%s.log", APP_NAME, APP_NAME)
	LogFolder          = fmt.Sprintf("/var/log/%s/", APP_NAME)
	TelegramBotChannel int
	DbPath             = "./data"
)
