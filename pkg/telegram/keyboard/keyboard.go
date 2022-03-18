package keyboard

import (
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/config"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard tg.ReplyKeyboardMarkup
type InKeyboard tg.InlineKeyboardMarkup

type Buttons struct {
	Tenant Tenant
	Admin  Admin
}

type Tenant struct {
	Keyboard Keyboard
	Water    []tg.InlineKeyboardButton
	Receipt  []tg.InlineKeyboardButton
}

type Admin struct {
	Keyboard Keyboard
	Rooms    InKeyboard
	Settings InKeyboard
}

func NewButtons() Buttons {
	cfg, _ := config.Init() //error no hand
	t := cfg.Text.Buttons.Tenant
	a := cfg.Text.Buttons.Admin
	return Buttons{
		Tenant: Tenant{
			Keyboard: Keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(t.Water1), tg.NewKeyboardButton(t.Receipt1)),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(t.Report1)))),

			Water: tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(t.Water.Hot_w2, t.Water.Hot_w2),
				tg.NewInlineKeyboardButtonData(t.Water.Cold_w2, t.Water.Cold_w2)),

			Receipt: tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(t.Receipt.Month2, t.Receipt.Month2),
				tg.NewInlineKeyboardButtonData(t.Receipt.Amount2, t.Receipt.Amount2),
				tg.NewInlineKeyboardButtonData(t.Receipt.Receipt2, t.Receipt.Receipt2),
			),
		},
		Admin: Admin{
			Keyboard: Keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(a.Rooms1), tg.NewKeyboardButton(a.Settings1)))),

			Rooms: InKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(a.Rooms.R2, a.Rooms.R2),
					tg.NewInlineKeyboardButtonData("366", "366")))),

			Settings: InKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(a.Settings.Edit2, a.Settings.Edit2)),
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(a.Settings.Contacts2, a.Settings.Contacts2)),
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(a.Settings.Reminder2, a.Settings.Reminder2)))),
		},
	}
}
