package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TenantStart struct{ CommandResponser }

func (r *TenantStart) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Response.Start)
	msg.ReplyMarkup = tgBot.But.Tenant.Keyboard
	_, err := tgBot.API.Send(msg)
	return err
}

type TenantCancel struct{ CommandResponser }

func (r *TenantCancel) Action(u *tg.Update) error {
	tgBot.State.Erase()
	return tgBot.API.SendText(u, "Операци отменена")
}

type TenantUnknownCmd struct{ CommandResponser }

func (r *TenantUnknownCmd) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_cmd)
}

type TenantHi struct{ MessageResponser }

func (r *TenantHi) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, fmt.Sprintf("(Tenant) Hello, %s!", u.Message.From.FirstName))
}

type TenantUnknownMes struct{ MessageResponser }

func (r *TenantUnknownMes) Action(u *tg.Update) error {
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

type Water1 struct{ ButtonResponser }

func (r *Water1) ShowButtons(u *tg.Update) error {
	tgBot.State.Erase()
	msg := tg.NewMessage(u.Message.Chat.ID, "Нажмите на нужный счёт за вводу и введите его значение")
	msg.ReplyMarkup = tgBot.But.Tenant.Water
	_, err := tgBot.API.Send(msg)
	return err
}

type Receipt1 struct{ ButtonResponser }

func (r *Receipt1) ShowButtons(u *tg.Update) error {
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
			inlineButtons = append(inlineButtons, tgBot.But.Tenant.Receipt[i])
		}
	}

	msg.ReplyMarkup = InKeyboard(tg.NewInlineKeyboardMarkup(inlineButtons))
	_, err := tgBot.API.Send(msg)
	return err
}

type Report1 struct{ ButtonResponser }

func (r *Report1) ShowButtons(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, "По техническим вопросам пишите в личные сообщения @MrDjeb")
	_, err := tgBot.API.Send(msg)
	return err
}

type Cold_w2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Cold_w2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите показания с счётчика холодной воды. К примеру: 34,56")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите показания с счётчика холодной воды. К примеру: 34,56"); err != nil {
		return err
	}
	tgBot.State.Erase()
	tgBot.State.TenantCold_w2 = true
	return nil
}

type Hot_w2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Hot_w2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите показания с счётчика горячей воды. К примеру: 34,56")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите показания с счётчика горячей воды. К примеру: 34,56"); err != nil {
		return err
	}
	tgBot.State.Erase()
	tgBot.State.TenantHot_w2 = true
	return nil
}

type Month2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Month2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите номер месяца -- число от 1 до 12.")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите номер месяца -- число от 1 до 12."); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[0] = 1
	case tgBot.State.TenantPayment[0] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[0] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[0] = 1
	}

	return nil
}

type Amount2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Amount2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500"); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[1] = 1
	case tgBot.State.TenantPayment[1] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[1] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[1] = 1
	}
	return nil
}

type Receipt2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Receipt2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Пришлите скрин квитанции.")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Пришлите скрин квитанции."); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[2] = 1
	case tgBot.State.TenantPayment[2] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[2] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[2] = 1
	}
	return nil
}
