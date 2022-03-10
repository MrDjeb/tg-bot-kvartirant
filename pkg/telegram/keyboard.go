package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdStart  = "start"
	cmdCancel = "cancel"
	msHi      = "Hi"
)

type keyboard tg.ReplyKeyboardMarkup
type inKeyboard tg.InlineKeyboardMarkup

type Buttons struct {
	Tenant
	Admin keyboard
}

type Tenant struct {
	keyboard keyboard
	Water    inKeyboard
	Receipt  []tg.InlineKeyboardButton
}

func (b *Bot) butInit() {
	b.But = Buttons{
		Tenant: Tenant{
			keyboard: keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Water1)), tg.NewKeyboardButton(string(b.Text.Receipt1))),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Report1))))),

			Water: inKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(string(b.Text.Water.Hot_w2), string(b.Text.Water.Hot_w2)), tg.NewInlineKeyboardButtonData(string(b.Text.Water.Cold_w2), string(b.Text.Water.Cold_w2))))),

			Receipt: tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(string(b.Text.Receipt.Add_month2), string(b.Text.Receipt.Add_month2)),
				tg.NewInlineKeyboardButtonData(string(b.Text.Receipt.Add_amount2), string(b.Text.Receipt.Add_amount2)),
				tg.NewInlineKeyboardButtonData(string(b.Text.Receipt.Add_receipt2), string(b.Text.Receipt.Add_receipt2)),
			),
		},
		Admin: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.Text.Admin.Rooms1))))),
	}
}

func (b *Bot) TenantHandlerClb(update *tg.Update) error {
	switch update.CallbackQuery.Data {
	case string(b.Text.Water.Cold_w2):
		return b.TenantCold_w2Clb(update)
	case string(b.Text.Water.Hot_w2):
		return b.TenantHot_w2Clb(update)
	case string(b.Text.Receipt.Add_amount2):
		return b.TenantAdd_amount2Clb(update)
	case string(b.Text.Receipt.Add_month2):
		return b.TenantAdd_month2Clb(update)
	case string(b.Text.Receipt.Add_receipt2):
		return b.TenantAdd_receipt2Clb(update)
	default:
		return b.handleSendText(update.CallbackQuery.From.ID, b.Text.Response.Unknown_ms)
	}
}

func (b *Bot) AdminHandlerClb(update *tg.Update) error {
	switch update.CallbackQuery.Data {
	case string(b.Text.Admin.Rooms1):
		return b.AdminUpClb(update)
	default:
		return b.handleSendText(update.CallbackQuery.From.ID, b.Text.Response.Unknown_ms)
	}
}

func (b *Bot) TenantHandlerCmd(message *tg.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.TenantStartCmd(message)
	case cmdCancel:
		return b.TenantCancelCmd(message)
	default:
		return b.handleSendText(message.From.ID, b.Text.Response.Unknown_cmd)
	}
}

func (b *Bot) AdminHandlerCmd(message *tg.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.AdminStartCmd(message)
	default:
		return b.handleSendText(message.From.ID, b.Text.Response.Unknown_cmd)
	}
}

func (b *Bot) TenantHandlerPh(message *tg.Message) error {
	if b.State.TenantPayment[2] == 1 {
		return b.TenantAdd_receipt2Inp(message)
	} else if b.State.TenantHot_w2 || b.State.TenantCold_w2 || b.State.TenantPayment[0] == 1 || b.State.TenantPayment[1] == 1 {
		return b.handleSendText(message.From.ID, "Сейчас мне не нужно фото.")
	} else {
		return b.handleSendText(message.From.ID, b.Text.Response.Unknown_ms)
	}
}

func (b *Bot) TenantHandlerMs(message *tg.Message) error {
	switch message.Text {
	case msHi:
		return b.TenantHiMs(message)
	case string(b.Text.Water1):
		return b.TenantWater1Ms(message)
	case string(b.Text.Receipt1):
		return b.TenantReceipt1Ms(message)
	case string(b.Text.Report1):
		return b.TenantReport1Ms(message)
	default:
		switch {
		case b.State.TenantHot_w2:
			return b.TenantHot_w2Inp(message)
		case b.State.TenantCold_w2:
			return b.TenantCold_w2Inp(message)
		case b.State.TenantPayment[0] == 1:
			return b.TenantAdd_month2Inp(message)
		case b.State.TenantPayment[1] == 1:
			return b.TenantAdd_amount2Inp(message)
		case b.State.TenantPayment[2] == 1:
			return b.handleSendText(message.From.ID, "Пришлите фото.")
		default:
			return b.TenantUnknownMs(message)
		}
	}
}

func (b *Bot) AdminHandlerMs(message *tg.Message) error {
	switch message.Text {
	case msHi:
		return b.AdminHiMs(message)
	default:
		return b.handleSendText(message.From.ID, b.Text.Response.Unknown_ms)
	}
}
