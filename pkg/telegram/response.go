package telegram

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownStart struct{ CommandResponser }

func (r *UnknownStart) Action(u *tg.Update) error {
	byteToken, err := base64.StdEncoding.DecodeString(u.Message.CommandArguments())
	if err != nil {
		return err
	}
	token := string(byteToken)
	idAdmin, err := strconv.ParseInt(token[32:], 10, 64)
	if err != nil {
		return err
	}

	st, ok := tgBot.State.Get(cache.KeyT(idAdmin))
	if !ok {
		return nil
	}
	d := st.Data.(cache.AdminData)
	number, ok := d.AddingRooms[token]
	if !ok {
		return tgBot.API.SendText(u, "Ссылка не валидная или её срок годности истёк.")
	}
	delete(d.AddingRooms, token)
	tgBot.State.Put(cache.KeyT(idAdmin), cache.State{Data: d})

	room := database.Room{
		IdTgAdmin:  database.TelegramID(idAdmin),
		IdTgTenant: database.TelegramID(u.FromChat().ID),
		Number:     database.Number(number),
	}
	if err := tgBot.DB.Room.Insert(room); err != nil {
		return err
	}

	tenant := database.Tenant{IdTg: database.TelegramID(u.FromChat().ID)}
	if err := tgBot.DB.Tenant.Insert(tenant); err != nil {
		return err
	}

	return tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.Cmd[tgBot.Text.CommonCommand.Start].Action(u)
}

type UnknownUnknownCmd struct{ CommandResponser }

func (r *UnknownUnknownCmd) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, "Обратитесь к администратору для регистрации. Вы не авторезированный пользователь.")
}

/////////////////////////////////////////

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

/////////////////////////////////////////////////////////////////////

type AdminStart struct {
	CommandResponser
	But keyboard.Keyboard
}

func (r *AdminStart) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Response.Start)
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type AdminCancel struct{ CommandResponser }

func (r *AdminCancel) Action(u *tg.Update) error {
	tgBot.State.Del(cache.KeyT(u.FromChat().ID))
	return tgBot.API.SendText(u, "Операци отменена")
}

type AdminUnknownCmd struct{ CommandResponser }

func (r *AdminUnknownCmd) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_cmd)
}

type AdminHi struct{ MessageResponser }

func (r *AdminHi) Action(u *tg.Update) error {
	return tgBot.API.SendText(u, fmt.Sprintf("(Admin) Hello, %s!", u.Message.From.FirstName))
}

type AdminUnknownMes struct{ MessageResponser }

func (r *AdminUnknownMes) Action(u *tg.Update) error {
	/*st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		if d := st.Data.(cache.AdminData); d.Payment[0] || d.Payment[1] || d.Payment[2] || d.Score[0] || d.Score[1] {
			return tgBot.API.SendText(u, "Продолжите внесение данных.")
		}
	}*/
	return tgBot.API.SendText(u, tgBot.Text.Unknown_ms)

}

type Rooms1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Rooms1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Выбирите комнату")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Settings1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Settings1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Выбирите настройку")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Edit2 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *Edit2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Скопируйте ссылку и отошлите её")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Скопируйте ссылку и отошлите её"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Edit2) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Выбирите")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Reminder2 struct{ InbuttonResponser }

func (r *Reminder2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Уведомления отправлены")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Уведомления отправлены"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Reminder2) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Sended notification")
	_, err := tgBot.API.Send(msg)
	return err
}

type Contacts2 struct{ InputResponser }

func (r *Contacts2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите ник в телеграмма")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите ник в телеграмма"); err != nil {
		return err
	}

	return nil
}

type AddRoom3 struct{ InputResponser }

func (r *AddRoom3) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите номер комнаты, который хотите добавить")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите номер комнаты, который хотите добавить"); err != nil {
		return err
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Admin.Settings.Edit.AddRoom3, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Admin.Settings.Edit.AddRoom3,
			Data: cache.AdminData{AddingRooms: make(map[string]string)}})
	}

	return nil
}

type RemoveRoom3 struct{ InbuttonResponser }

func (r *RemoveRoom3) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите номер комнаты, который хотите добавить")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите номер комнаты, который хотите добавить"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *RemoveRoom3) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Выберите комнату комнаты")
	_, err := tgBot.API.Send(msg)
	return err
}

type ShowScorer33 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ShowScorer33) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Последние 3 показания счётчика")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Последние три показания счётчика"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorer33) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Счёт зав воду 1, 2 ")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type ShowScorer14 struct {
	InbuttonResponser
}

func (r *ShowScorer14) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Последние показание счётчика")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Последние показание счётчика"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorer14) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Счёт зав воду")
	_, err := tgBot.API.Send(msg)
	return err
}

type ShowScorerN4 struct {
	InbuttonResponser
}

func (r *ShowScorerN4) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Показания счётчиков")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Показания счётчиков"); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorerN4) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Счёт зав воду все")
	_, err := tgBot.API.Send(msg)
	return err
}
