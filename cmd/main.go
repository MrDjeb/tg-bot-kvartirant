package main

import (
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
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

	db, err := database.Init()
	cherr(err)
	defer db.Scorer.DB.Close()

	cach, err := cache.NewCache(60 * 60 * 60)
	cherr(err)
	defer cach.Destroy()

	botAPI, err := tg.NewBotAPI(cfg.TgToken)
	cherr(err)

	botAPI.Debug = true

	tgBot := telegram.NewBot(botAPI, cfg.Text, db, cach)
	cherr(tgBot.Start())
}
