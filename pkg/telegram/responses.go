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

func (b *Bot) TenantAdd_month2Clb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Введите номер месяца -- число от 1 до 12.")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите номер месяца -- число от 1 до 12.")
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(b.State.TenantPayment[0]) && !tbool(b.State.TenantPayment[1]) && !tbool(b.State.TenantPayment[2]):
		b.State.Erase()
		b.State.TenantPayment[0] = 1
	case b.State.TenantPayment[0] == 2:
		b.State.Erase()
		b.State.TenantPayment[0] = 1
	default:
		b.State.CleanProcess()
		b.State.TenantPayment[0] = 1
	}

	return nil
}

func (b *Bot) TenantAdd_amount2Clb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500")
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(b.State.TenantPayment[0]) && !tbool(b.State.TenantPayment[1]) && !tbool(b.State.TenantPayment[2]):
		b.State.Erase()
		b.State.TenantPayment[1] = 1
	case b.State.TenantPayment[1] == 2:
		b.State.Erase()
		b.State.TenantPayment[1] = 1
	default:
		b.State.CleanProcess()
		b.State.TenantPayment[1] = 1
	}
	return nil
}

func (b *Bot) TenantAdd_receipt2Clb(update *tg.Update) error {
	callback := tg.NewCallback(update.CallbackQuery.ID, "Пришлите скрин квитанции.")
	if _, err := b.Api.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пришлите скрин квитанции.")
	if _, err := b.Api.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(b.State.TenantPayment[0]) && !tbool(b.State.TenantPayment[1]) && !tbool(b.State.TenantPayment[2]):
		b.State.Erase()
		b.State.TenantPayment[2] = 1
	case b.State.TenantPayment[2] == 2:
		b.State.Erase()
		b.State.TenantPayment[2] = 1
	default:
		b.State.CleanProcess()
		b.State.TenantPayment[2] = 1
	}
	return nil
}

func (b *Bot) TenantStartCmd(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.Text.Response.Start)
	msg.ReplyMarkup = b.But.Tenant.keyboard
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantCancelCmd(message *tg.Message) error {
	b.State.Erase()
	msg := tg.NewMessage(message.Chat.ID, "Операция отменена.")
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
	b.State.Erase()
	msg := tg.NewMessage(message.Chat.ID, "Нажмите на нужный счёт за вводу и введите его значение")
	msg.ReplyMarkup = b.But.Water
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantReceipt1Ms(message *tg.Message) error {
	var msg tg.MessageConfig
	sumP := (b.State.TenantPayment[0] + b.State.TenantPayment[1] + b.State.TenantPayment[2]) / 2
	switch sumP {
	case 0: // 0
		b.State.Erase()
		msg = tg.NewMessage(message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести первым.")
	case 1: // 2
		msg = tg.NewMessage(message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести следующим.")
	case 2: // 4
		msg = tg.NewMessage(message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВнесите последний параметр.")
	}

	var inlineButtons []tg.InlineKeyboardButton
	for i := range b.State.TenantPayment {
		if !tbool(b.State.TenantPayment[i]) {
			inlineButtons = append(inlineButtons, b.But.Receipt[i])
		}
	}

	msg.ReplyMarkup = inKeyboard(tg.NewInlineKeyboardMarkup(inlineButtons))
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantReport1Ms(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, "По техническим вопросам пишите в личные сообщения @MrDjeb")
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) TenantUnknownMs(message *tg.Message) error {
	sumP := b.State.TenantPayment[0] + b.State.TenantPayment[1] + b.State.TenantPayment[2]
	if sumP > 0 {
		msg := tg.NewMessage(message.Chat.ID, "Продолжите внесение данных.")
		_, err := b.Api.Send(msg)
		return err
	} else {
		msg := tg.NewMessage(message.Chat.ID, b.Text.Unknown_ms)
		_, err := b.Api.Send(msg)
		return err
	}
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
	msg.ReplyMarkup = b.But.Admin
	_, err := b.Api.Send(msg)
	return err
}

func (b *Bot) AdminHiMs(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, fmt.Sprintf("(Admin) Hello, %s!", message.From.FirstName))
	_, err := b.Api.Send(msg)
	return err
}
