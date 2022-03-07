package telegram

import (
	"errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleBack(update *tg.Update) error {
	flagT, flagA := b.DB.Tenant.IsExist(update.CallbackQuery.From.ID), b.DB.Admin.IsExist(update.CallbackQuery.From.ID)
	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: userID in double table")
	case flagT:
		return b.TenantHandlerClb(update)
	case flagA:
		return b.AdminHandlerClb(update)
	default:
		return nil
	}
}

func (b *Bot) handleCmd(message *tg.Message) error {
	flagT, flagA := b.DB.Tenant.IsExist(message.From.ID), b.DB.Admin.IsExist(message.From.ID)
	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: userID in double table")
	case flagT:
		return b.TenantHandlerCmd(message)
	case flagA:
		return b.AdminHandlerCmd(message)
	default:
		return nil
	}
}

func (b *Bot) handleMs(message *tg.Message) error {
	flagT, flagA := b.DB.Tenant.IsExist(message.From.ID), b.DB.Admin.IsExist(message.From.ID)
	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: UserID in double table")
	case flagT:
		return b.TenantHandlerMs(message)
	case flagA:
		return b.AdminHandlerMs(message)
	default:
		return nil
	}
}

func (b *Bot) handleSendText(message *tg.Message, text string) error {
	msg := tg.NewMessage(message.Chat.ID, text)
	_, err := b.Api.Send(msg)
	return err
}
