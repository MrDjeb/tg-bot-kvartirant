package telegram

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgBot *Bot

type User struct {
	Handler
}

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

func (a API) DelMes(u *tg.Update) error {
	delete := tg.NewDeleteMessage(u.FromChat().ID, u.CallbackQuery.Message.MessageID)
	if _, err := tgBot.API.Request(delete); err != nil {
		return err
	}
	return nil
}

func (a API) AnsCallback(u *tg.Update, text string) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, text)
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	return nil
}

type logBot struct {
	std *log.Logger
}

func (l *logBot) Println(v ...interface{}) { l.std.Output(2, DecodeUTF16(fmt.Sprintln(v...))) }

func (l *logBot) Printf(format string, v ...interface{}) {
	l.std.Output(2, DecodeUTF16(fmt.Sprintf(format, v...)))
}

func (b *Bot) Start() error {
	tg.SetLogger(&logBot{log.New(os.Stderr, "[API] ", log.LstdFlags|log.Lmsgprefix)})

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)
	defer b.API.StopReceivingUpdates()

	for u := range updates {
		user, err := b.FromWhom(&u)
		if err != nil {
			return err
		}

	out:
		switch {
		case u.CallbackQuery != nil:
			if err := user.Callback(&u); err != nil {
				return err
			}
			break out
		case u.Message.IsCommand():
			if err := user.Command(&u); err != nil {
				return err
			}
			break out
		case u.Message.Photo != nil:
			if err := user.Photo(&u); err != nil {
				return err
			}
			break out
		case u.Message.Text != "":
			if err := user.Message(&u); err != nil {
				return err
			}
			break out
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

func DecodeUTF16(str string) string {
	str = strings.Replace(str, "\"", "＂", -1)
	str = strings.Replace(str, "\n", "", -1)

	s, err := strconv.Unquote("\"" + str + "\"")
	if err != nil {
		log.Println(err)
	}

	return s
}
