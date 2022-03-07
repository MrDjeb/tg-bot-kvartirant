package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) TenantCold_w2Clb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Введите показания с счётчика холодной воды. К примеру: 34,56")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите показания с счётчика холодной воды. К примеру: 34,56")
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}
	b.State.Erase()
	b.State.TenantCold_w2 = true
	return nil
}

func (b *Bot) TenantHot_w2Clb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Введите показания с счётчика горячей воды. К примеру: 34,56")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите показания с счётчика горячей воды. К примеру: 34,56")
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}
	b.State.Erase()
	b.State.TenantHot_w2 = true
	return nil
}

func (b *Bot) TenantStartCmd(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.Text.Response.Start)
	msg.ReplyMarkup = b.But.Tenant.keyboard
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantHiMs(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, fmt.Sprintf("(Tenant) Hello, %s!", message.From.FirstName))
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantWater1Ms(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, "Нажмите на нужный счёт за вводу и введите его значение")
	msg.ReplyMarkup = b.But.Water.keyboard
	_, err := b.Api.Send(msg)
	return err
}

////////////////
func (b *Bot) AdminUpClb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Adminnn!")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) AdminStartCmd(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.Text.Response.Start)
	msg.ReplyMarkup = b.But.Admin.keyboard
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) AdminHiMs(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, fmt.Sprintf("(Admin) Hello, %s!", message.From.FirstName))
	_, err := b.Api.Send(msg)
	return err
}
