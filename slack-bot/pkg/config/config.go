package config

import (
	"os"

	"github.com/joho/godotenv"

	"slack-bot/pkg/utils"
)

var (
	SLACK_BOT_TOKEN string
	SLACK_APP_TOKEN string
)

func init() {
	loadEnv()
	setEnv()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Fatalln("Failed to load .env file")
	}
}

func setEnv() {
	SLACK_APP_TOKEN = os.Getenv("SLACK_APP_TOKEN")
	SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")

	if SLACK_APP_TOKEN == "" || SLACK_BOT_TOKEN == "" {
		utils.ErrorLogger.Fatalf("[ERROR]: SLACK_BOT_TOKEN or SLACK_APP_TOKEN are not set in .env file")
	}
}
