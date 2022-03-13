package telegram

import (
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgBot *Bot

type Bot struct {
	API   API
	Text  config.Text
	State State
	But   Buttons
	DB    database.Tables
	Handler
	Tenant User
	Admin  User
}

func NewBot(api *tg.BotAPI, text config.Text, db database.Tables) *Bot {
	b := &Bot{
		API:    API{api},
		Text:   text,
		State:  State{},
		But:    Buttons{},
		DB:     db,
		Tenant: User{},
		Admin:  User{},
	}
	b.State.Erase()
	return b
}

func (b *Bot) Init() {
	b.But = NewButtons()
	b.Tenant = User{NewTenantHandler()}
	b.Admin = User{NewTenantHandler()}
}

type API struct {
	*tg.BotAPI
}

func (a API) SendText(u *tg.Update, text string) error {
	msg := tg.NewMessage(u.SentFrom().ID, text)
	_, err := a.Send(msg)
	return err
}

type User struct {
	Handler
}

func (b *Bot) Start() error {
	/*if err := b.DB.Tenant.Insert(database.Tenant{IdTg: 410345981}); err != nil {
		return err
	}*/

	tgBot = b
	tgBot.Init()

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for u := range updates {

		switch {
		case u.CallbackQuery != nil:
			if err := b.Callback(&u); err != nil {
				return err
			}
			continue
		case u.Message.IsCommand():
			if err := b.Command(&u); err != nil {
				return err
			}
			continue
		case u.Message.Photo != nil:
			if err := b.Photo(&u); err != nil {
				return err
			}
			continue
		case u.Message.Text != "":
			if err := b.Message(&u); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func (b *Bot) FromWhom(u *tg.Update) (bool, bool, error) {
	flagT, err := b.DB.Tenant.IsExist(u.SentFrom().ID)
	if err != nil {
		return false, false, nil
	}
	flagA, err := b.DB.Admin.IsExist(u.SentFrom().ID)
	if err != nil {
		return false, false, err
	}
	return flagT, flagA, nil
}
