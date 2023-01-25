package main

import (
	"os"
	"awesomeProject/telegram"
	"log"
)

func main() {
	bot, err := telegram.CreateTelegramBot(os.Getenv("tokenApi"))

	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
