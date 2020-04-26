package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ENV                string
	API_TOKEN          string
	APP_NAME           = "assistme"
	LogFolder          = fmt.Sprintf("/var/log/%s", APP_NAME)
	LogPath            = fmt.Sprintf("%s/%s.log", LogFolder, APP_NAME)
	TelegramBotChannel int
	DbPath             string
)

const (
	// key delimiter
	KeyDelim = ":"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Could not load dotenv file: %+v\n", err))
	}
	ENV = os.Getenv("ENV")
	cwd, err := os.Getwd()
	if err != nil {
		panic("Could not get working directory")
	}
	if ENV == "DEV" {
		LogFolder = cwd
		LogPath = fmt.Sprintf("%s/%s.log", LogFolder, APP_NAME)
	}
	DbPath = fmt.Sprintf("%s/data/", cwd)
	API_TOKEN = os.Getenv("TELEGRAM_APITOKEN")
	TelegramBotChannel, err = strconv.Atoi(os.Getenv("BOT_CHANNEL_ID"))
	if err != nil {
		panic("Could not parse Bot channel id")
	}

}
