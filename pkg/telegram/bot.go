package telegram

import (
	"database/sql"
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var DBScorer database.DBScorer

type Bot struct {
	bot     *tg.BotAPI
	text    config.Text
	buttons Buttons
}

func NewBot(bot *tg.BotAPI, text config.Text) *Bot {
	return &Bot{
		bot:     bot,
		text:    text,
		buttons: Buttons{},
	}
}

func (b *Bot) Start() error {
	b.buttons = b.butInit()
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	//InitDB()

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				return err
			}

			continue
		}

		if update.Message.Text != "" {
			if err := b.handleMessage(update.Message); err != nil {
				return err

			}
			continue
		}

		//_, err := DBScorer.Insert(database.Scorer{update.Message.From.ID, 0, 0, time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)})
		//cherr(err)

	}
	return nil
}

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func InitDB() {
	db, err := sql.Open("sqlite3", "bot.db")
	cherr(err)
	DBScorer = database.DBScorer{DB: db}
}
