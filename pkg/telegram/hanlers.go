package telegram

import (
	"errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleBack(update *tg.Update) error {
	flagT, err := b.DB.Tenant.IsExist(update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}
	flagA, err := b.DB.Admin.IsExist(update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}

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
	flagT, err := b.DB.Tenant.IsExist(message.From.ID)
	if err != nil {
		return err
	}
	flagA, err := b.DB.Admin.IsExist(message.From.ID)
	if err != nil {
		return err
	}

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

func (b *Bot) handlePh(message *tg.Message) error {
	flagT, err := b.DB.Tenant.IsExist(message.From.ID)
	if err != nil {
		return err
	}
	flagA, err := b.DB.Admin.IsExist(message.From.ID)
	if err != nil {
		return err
	}

	switch {
	case flagT && flagA:
		return errors.New("incorrect database data: UserID in double table")
	case flagT:
		return b.TenantHandlerPh(message)
	default:
		return nil
	}
}

func (b *Bot) handleMs(message *tg.Message) error {
	flagT, err := b.DB.Tenant.IsExist(message.From.ID)
	if err != nil {
		return err
	}
	flagA, err := b.DB.Admin.IsExist(message.From.ID)
	if err != nil {
		return err
	}

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

func (b *Bot) handleSendText(id int64, text string) error {
	msg := tg.NewMessage(id, text)
	_, err := b.Api.Send(msg)
	return err
}
