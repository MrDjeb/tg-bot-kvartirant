package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdStart = "start"
	msHi     = "Hi"
)

type keyboard tg.ReplyKeyboardMarkup
type inKeyboard tg.InlineKeyboardMarkup

type Buttons struct {
	Tenant
	Admin
}

type Tenant struct {
	keyboard keyboard
	Water
	Receipt
}

type Admin struct {
	keyboard keyboard
}

type Water struct {
	keyboard inKeyboard
}

type Receipt struct {
	add_month2   keyboard
	add_amount2  keyboard
	add_receipt2 keyboard
}

func (b *Bot) butInit() {
	b.But = Buttons{
		Tenant: Tenant{
			keyboard: keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Water1)), tg.NewKeyboardButton(string(b.Text.Receipt1))),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Report1))))),

			Water: Water{
				inKeyboard(tg.NewInlineKeyboardMarkup(
					tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(string(b.Text.Water.Hot_w2), string(b.Text.Water.Hot_w2)), tg.NewInlineKeyboardButtonData(string(b.Text.Water.Cold_w2), string(b.Text.Water.Cold_w2))))),
			},

			Receipt: Receipt{
				add_month2:   keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Receipt.Add_month2))))),
				add_amount2:  keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Receipt.Add_amount2))))),
				add_receipt2: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Receipt.Add_receipt2))))),
			},
		},
		Admin: Admin{
			keyboard: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Admin.Rooms1))))),
		},
	}
}

func (b *Bot) TenantHandlerClb(update *tg.Update) error {
	switch update.CallbackQuery.Data {
	case string(b.Text.Water.Cold_w2):
		return b.TenantCold_w2Clb(update)
	case string(b.Text.Water.Hot_w2):
		return b.TenantHot_w2Clb(update)
	default:
		return b.handleSendText(update.Message, b.Text.Response.Unknown_ms)
	}
}

func (b *Bot) AdminHandlerClb(update *tg.Update) error {
	switch update.CallbackQuery.Data {
	case string(b.Text.Admin.Rooms1):
		return b.AdminUpClb(update)
	default:
		return b.handleSendText(update.Message, b.Text.Response.Unknown_ms)
	}
}

func (b *Bot) TenantHandlerCmd(message *tg.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.TenantStartCmd(message)
	default:
		return b.handleSendText(message, b.Text.Response.Unknown_cmd)
	}
}

func (b *Bot) AdminHandlerCmd(message *tg.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.AdminStartCmd(message)
	default:
		return b.handleSendText(message, b.Text.Response.Unknown_cmd)
	}
}

func (b *Bot) TenantHandlerMs(message *tg.Message) error {
	switch {
	case b.State.TenantHot_w2:
		return b.TenantHot_w2Inp(message)
	case b.State.TenantCold_w2:
		return b.TenantCold_w2Inp(message)
	default:
		switch message.Text {
		case msHi:
			return b.TenantHiMs(message)
		case string(b.Text.Water1):
			return b.TenantWater1Ms(message)
		default:
			return b.handleSendText(message, b.Text.Response.Unknown_ms)
		}
	}
}

func (b *Bot) AdminHandlerMs(message *tg.Message) error {
	switch message.Text {
	case msHi:
		return b.AdminHiMs(message)
	default:
		return b.handleSendText(message, b.Text.Response.Unknown_ms)
	}
}
