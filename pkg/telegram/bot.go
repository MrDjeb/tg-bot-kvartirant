package telegram

import (
	"errors"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgBot *Bot

type Bot struct {
	API   API
	Text  config.Text
	State *cache.Cache
	DB    database.Tables
	Handler
	Tenant  User
	Admin   User
	Unknown User
}

func NewBot(api *tg.BotAPI, text config.Text, db database.Tables, cach *cache.Cache) *Bot {
	b := &Bot{
		API:     API{api},
		Text:    text,
		State:   cach,
		DB:      db,
		Tenant:  User{},
		Admin:   User{},
		Unknown: User{},
	}
	tgBot = b
	b.Tenant = NewUser(&TenantHandler{})
	b.Admin = NewUser(&AdminHandler{})
	b.Unknown = NewUser(&UnknownHandler{})
	return b
}

type API struct {
	*tg.BotAPI
}

func (a API) SendText(u *tg.Update, text string) error {
	msg := tg.NewMessage(u.FromChat().ID, text)
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

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)
	defer b.API.StopReceivingUpdates()

	for u := range updates {
		user, err := b.FromWhom(&u)
		if err != nil {
			return err
		}

		switch {
		case u.CallbackQuery != nil:
			if err := user.Callback(&u); err != nil {
				return err
			}
		case u.Message.IsCommand():
			if err := user.Command(&u); err != nil {
				return err
			}
		case u.Message.Photo != nil:
			if err := user.Photo(&u); err != nil {
				return err
			}
		case u.Message.Text != "":
			if err := user.Message(&u); err != nil {
				return err
			}
		}
	}
	return errors.New("there is no handler for this update")
}

func (b *Bot) FromWhom(u *tg.Update) (User, error) {
	flagT, err := b.DB.Tenant.IsExist(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return User{}, err
	}

	flagA, err := b.DB.Admin.IsExist(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return User{}, err
	}

	switch {
	case flagA && flagT:
		return User{}, errors.New("intersection database IdTelegram in data")
	case flagT:
		return tgBot.Tenant, err
	case flagA:
		return tgBot.Admin, err
	default:
		return tgBot.Unknown, err
	}
}
