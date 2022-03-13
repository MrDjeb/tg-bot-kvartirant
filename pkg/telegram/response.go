package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TenantResponser struct {
	Cmd TenantResponseCommand
	Ms  TenantResponseMessage
}

func (i *TenantResponser) New() {
	i.Cmd = TenantResponseCommand{}
	i.Ms = TenantResponseMessage{}
}

type TenantResponseCommand struct{}

func (r *TenantResponseCommand) Start(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Response.Start)
	msg.ReplyMarkup = tgBot.But.Tenant.keyboard
	_, err := tgBot.API.Send(msg)
	return err
}
func (r *TenantResponseCommand) Cancel(u *tg.Update) error {
	tgBot.State.Erase()
	return tgBot.API.SendText(u, "Операци отменена")
}
func (r *TenantResponseCommand) Unknown(u *tg.Update) error {
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_cmd)
}

type TenantResponseMessage struct{}

func (r *TenantResponseMessage) Hi(u *tg.Update) error {
	return tgBot.API.SendText(u, fmt.Sprintf("(Tenant) Hello, %s!", u.Message.From.FirstName))
}

func (r *TenantResponseMessage) Water1(u *tg.Update) error {
	tgBot.State.Erase()
	msg := tg.NewMessage(u.Message.Chat.ID, "Нажмите на нужный счёт за вводу и введите его значение")
	msg.ReplyMarkup = tgBot.But.Water
	_, err := tgBot.API.Send(msg)
	return err
}

func (r *TenantResponseMessage) Receipt1(u *tg.Update) error {
	var msg tg.MessageConfig
	sumP := (tgBot.State.TenantPayment[0] + tgBot.State.TenantPayment[1] + tgBot.State.TenantPayment[2]) / 2
	switch sumP {
	case 0: // 0
		tgBot.State.Erase()
		msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести первым.")
	case 1: // 2
		msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести следующим.")
	case 2: // 4
		msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВнесите последний параметр.")
	}

	var inlineButtons []tg.InlineKeyboardButton
	for i := range tgBot.State.TenantPayment {
		if !tbool(tgBot.State.TenantPayment[i]) {
			inlineButtons = append(inlineButtons, tgBot.But.Receipt[i])
		}
	}

	msg.ReplyMarkup = inKeyboard(tg.NewInlineKeyboardMarkup(inlineButtons))
	_, err := tgBot.API.Send(msg)
	return err
}

func (r *TenantResponseMessage) Report1(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, "По техническим вопросам пишите в личные сообщения @MrDjeb")
	_, err := tgBot.API.Send(msg)
	return err
}

func (r *TenantResponseMessage) Unknown(u *tg.Update) error {
	sumP := tgBot.State.TenantPayment[0] + tgBot.State.TenantPayment[1] + tgBot.State.TenantPayment[2]
	if sumP > 0 {
		msg := tg.NewMessage(u.Message.Chat.ID, "Продолжите внесение данных.")
		_, err := tgBot.API.Send(msg)
		return err
	} else {
		msg := tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Unknown_ms)
		_, err := tgBot.API.Send(msg)
		return err
	}
}
