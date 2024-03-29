package main

import (
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[ERROR] ")

	cfg, err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Scorer.DB.Close()

	cach, err := cache.NewCache(60 * 60 * 24)
	if err != nil {
		log.Fatalln(err)
	}
	defer cach.Destroy()

	botAPI, err := tg.NewBotAPI(cfg.TgToken)
	if err != nil {
		log.Fatalln(err)
	}

	botAPI.Debug = true
	cache.Debug = true

	tgBot := telegram.NewBot(botAPI, cfg.Text, db, cach)
	if err := tgBot.Start(); err != nil {
		log.Fatalln(err)
	}
}
