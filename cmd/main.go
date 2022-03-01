package main

import (
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	cfg, err := config.Init()
	cherr(err)

	bot, err := tg.NewBotAPI(cfg.TgToken)
	cherr(err)

	bot.Debug = true

	tgBot := telegram.NewBot(bot)
	cherr(tgBot.Start())

}
