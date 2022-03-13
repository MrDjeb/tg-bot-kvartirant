package telegram

import (
	"errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Callback(u *tg.Update) error
	Command(u *tg.Update) error
	Photo(u *tg.Update) error
	Message(u *tg.Update) error
}

func (h *Bot) Callback(u *tg.Update) error {
	flagT, flagA, err := h.FromWhom(u)
	if err != nil {
		return err
	}

	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: userID in double table")
	case flagT:
		return h.Tenant.Callback(u)
	case flagA:
		return h.Admin.Callback(u)
	default:
		return nil
	}
}

func (h *Bot) Command(u *tg.Update) error {
	flagT, flagA, err := h.FromWhom(u)
	if err != nil {
		return err
	}

	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: userID in double table")
	case flagT:
		return h.Tenant.Command(u)
	case flagA:
		return h.Admin.Command(u)
	default:
		return nil
	}
}

func (h *Bot) Photo(u *tg.Update) error {
	flagT, flagA, err := h.FromWhom(u)
	if err != nil {
		return err
	}

	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: UserID in double table")
	case flagT:
		return h.Tenant.Photo(u)
	case flagA:
		return h.Admin.Photo(u)
	default:
		return nil
	}
}

func (h *Bot) Message(u *tg.Update) error {
	flagT, flagA, err := h.FromWhom(u)
	if err != nil {
		return err
	}

	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: UserID in double table")
	case flagT:
		return h.Tenant.Message(u)
	case flagA:
		return h.Admin.Message(u)
	default:
		return nil
	}
}

type TenantHandler struct {
	Inp TenantInlineInput
	TenantResponser
}

func NewTenantHandler() *TenantHandler {
	var InlineInput TenantInlineInput
	InlineInput.New()
	var Responser TenantResponser
	Responser.New()
	return &TenantHandler{Inp: InlineInput, TenantResponser: Responser}
}

func (h *TenantHandler) Callback(u *tg.Update) error {
	switch u.CallbackQuery.Data {
	case string(tgBot.Text.Cold_w2):
		return h.Inp.Cold_w2.Callback(u)
	case string(tgBot.Text.Water.Hot_w2):
		return h.Inp.Hot_w2.Callback(u)
	case string(tgBot.Text.Receipt.Add_amount2):
		return h.Inp.Amount2.Callback(u)
	case string(tgBot.Text.Receipt.Add_month2):
		return h.Inp.Month2.Callback(u)
	case string(tgBot.Text.Receipt.Add_receipt2):
		return h.Inp.Receipt2.Callback(u)
	default:
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
	}
}
func (h *TenantHandler) Command(u *tg.Update) error {
	switch u.Message.Command() {
	case "start":
		return h.Cmd.Start(u)
	case "cancel":
		return h.Cmd.Cancel(u)
	default:
		return h.Cmd.Unknown(u)
	}
}
func (h *TenantHandler) Photo(u *tg.Update) error {
	if tgBot.State.TenantPayment[2] == 1 {
		return h.Inp.Receipt2.HandleInput(u)
	} else if tgBot.State.TenantHot_w2 || tgBot.State.TenantCold_w2 || tgBot.State.TenantPayment[0] == 1 || tgBot.State.TenantPayment[1] == 1 {
		return tgBot.API.SendText(u, "Сейчас мне не нужно фото.")
	} else {
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
	}
}
func (h *TenantHandler) Message(u *tg.Update) error {
	switch u.Message.Text {
	case "Hi":
		return h.Ms.Hi(u)
	case string(tgBot.Text.Water1):
		return h.Ms.Water1(u)
	case string(tgBot.Text.Receipt1):
		return h.Ms.Receipt1(u)
	case string(tgBot.Text.Report1):
		return h.Ms.Report1(u)
	default:
		switch {
		case tgBot.State.TenantHot_w2:
			return h.Inp.Hot_w2.HandleInput(u)
		case tgBot.State.TenantCold_w2:
			return h.Inp.Cold_w2.HandleInput(u)
		case tgBot.State.TenantPayment[0] == 1:
			return h.Inp.Month2.HandleInput(u)
		case tgBot.State.TenantPayment[1] == 1:
			return h.Inp.Amount2.HandleInput(u)
		case tgBot.State.TenantPayment[2] == 1:
			return tgBot.API.SendText(u, "Пришлите фото.")
		default:
			return h.Ms.Unknown(u)
		}
	}
}
