package main

import (
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	bot, err := tg.NewBotAPI("5150854501:AAHM8auF6KgpeHIbw2BHSVMJ5CRPshzYU5s")
	cherr(err)

	bot.Debug = true

	tgBot := telegram.NewBot(bot)
	cherr(tgBot.Start())

}
