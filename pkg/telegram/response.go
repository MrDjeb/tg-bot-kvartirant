package telegram

import (
	"fmt"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TenantStart struct {
	CommandResponser
	But keyboard.Keyboard
}

func (r *TenantStart) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Response.Start)
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type TenantCancel struct{ CommandResponser }

func (r *TenantCancel) Action(u *tg.Update) error {
	tgBot.State.Del(cache.KeyT(u.FromChat().ID))
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
	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		if d := st.Data.(cache.TenantData); d.Payment[0] || d.Payment[1] || d.Payment[2] || d.Score[0] || d.Score[1] {
			return tgBot.API.SendText(u, "Продолжите внесение данных.")
		}
	}
	return tgBot.API.SendText(u, tgBot.Text.Unknown_ms)

}

type Water1 struct {
	ButtonResponser
	But []tg.InlineKeyboardButton
}

func (r *Water1) Action(u *tg.Update) error {
	var msg tg.MessageConfig
	var inlineButtons []tg.InlineKeyboardButton
	fmt.Println("WAAAAAAAAAAAAAAAAATTTTTTEEEEEEERRRR!!!!!!!!!!")
	tgBot.State.Display()

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		switch tint(d.Score[0]) + tint(d.Score[1]) {
		case 0:
			msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех двух параметров.\nВыберете какой хотите внести первым.")
		case 1:
			msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех двух параметров.\nВнесите последний параметр.")
		}
		for i := range d.Score {
			if !d.Score[i] {
				inlineButtons = append(inlineButtons, r.But[i])
			}
		}
	} else {
		msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех двух параметров.\nВыберете какой хотите внести первым.")
		inlineButtons = r.But
	}

	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(inlineButtons)
	_, err := tgBot.API.Send(msg)
	return err
}

type Receipt1 struct {
	ButtonResponser
	But []tg.InlineKeyboardButton
}

func (r *Receipt1) Action(u *tg.Update) error {
	var msg tg.MessageConfig
	var inlineButtons []tg.InlineKeyboardButton

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		switch tint(d.Payment[0]) + tint(d.Payment[1]) + tint(d.Payment[2]) {
		case 0:
			msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести первым.")
		case 1:
			msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести следующим.")
		case 2:
			msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВнесите последний параметр.")
		}
		for i := range d.Payment {
			if !d.Payment[i] {
				inlineButtons = append(inlineButtons, r.But[i])
			}
		}
	} else {
		msg = tg.NewMessage(u.Message.Chat.ID, "Данные сохронятся после заполнения всех трёх параметров.\nВыберете какой хотите внести первым.")
		inlineButtons = r.But
	}

	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(inlineButtons)
	_, err := tgBot.API.Send(msg)
	return err
}

type Report1 struct{ ButtonResponser }

func (r *Report1) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, "По техническим вопросам пишите в личные сообщения @MrDjeb")
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

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Water.Hot_w2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Water.Hot_w2, Data: cache.TenantData{}})
	}

	return nil
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

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Water.Cold_w2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Water.Cold_w2, Data: cache.TenantData{}})
	}
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

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Month2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Month2, Data: cache.TenantData{}})
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

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Amount2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Amount2, Data: cache.TenantData{}})
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

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.TenantData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Receipt2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Tenant.Receipt.Receipt2, Data: cache.TenantData{}})
	}

	return nil
}
