package telegram

import (
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api   *tg.BotAPI
	Text  config.Text
	But   Buttons
	DB    database.Tables
	State State
}

func NewBot(api *tg.BotAPI, text config.Text, db database.Tables) *Bot {
	b := &Bot{
		Api:   api,
		Text:  text,
		But:   Buttons{},
		DB:    db,
		State: State{},
	}
	b.butInit()
	b.State.Erase()
	return b
}

func (b *Bot) Start() error {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.Api.GetUpdatesChan(u)

	if err := b.DB.Tenant.Migrate(); err != nil {
		return err
	}

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			if err := b.handleBack(&update); err != nil {
				return err
			}
			continue
		case update.Message.IsCommand():
			if err := b.handleCmd(update.Message); err != nil {
				return err
			}
			continue
		case update.Message.Text != "":
			if err := b.handleMs(update.Message); err != nil {
				return err
			}
			continue
		}

		//_, err := DBScorer.Insert(database.Scorer{update.Message.From.ID, 0, 0, time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)})
		//cherr(err)

	}
	return nil
}
