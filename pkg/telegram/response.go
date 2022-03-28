package telegram

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/olekukonko/tablewriter"
)

type UnknownStart struct{ CommandResponser }

func (r *UnknownStart) Action(u *tg.Update) error {
	if len(u.Message.CommandArguments()) < 32 {
		return tgBot.Unknown.Handler.(*UnknownHandler).HandlerResponse.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
	}

	byteToken, err := base64.StdEncoding.DecodeString(u.Message.CommandArguments())
	if err != nil { //error broke
		return tgBot.API.SendText(u, "Ссылка не валидная или её срок годности истёк.")
	}
	token := string(byteToken)
	idAdmin, err := strconv.ParseInt(token[32:], 10, 64)
	if err != nil { //error broke
		return tgBot.API.SendText(u, "Ссылка не валидная или её срок годности истёк.")
	}

	st, ok := tgBot.State.Get(cache.KeyT(idAdmin))
	if !ok {
		return tgBot.API.SendText(u, "Кэш с таким idAdmin пуст.")
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

	msg := tg.NewMessage(idAdmin, fmt.Sprintf("🔗 %s %s успешно привязался(-aсь) к комнате № 〈%s〉", u.SentFrom().FirstName, u.SentFrom().UserName, number))
	_, err = tgBot.API.Send(msg)
	if err != nil {
		return err
	}

	tenant := database.Tenant{IdTg: database.TelegramID(u.FromChat().ID)}
	if err := tgBot.DB.Tenant.Insert(tenant); err != nil {
		return err
	}

	return tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.Cmd[tgBot.Text.CommonCommand.Start].Action(u)
}

type GodMode struct{ CommandResponser }

func (r *GodMode) Action(u *tg.Update) error {
	if u.Message.CommandArguments() == tgBot.Text.GodModeKey {
		if err := tgBot.API.SendText(u, "Activated!"); err != nil {
			return err
		}
		if err := tgBot.DB.Admin.Insert(database.Admin{IdTgAdmin: database.TelegramID(u.FromChat().ID), Repairer: u.SentFrom().UserName}); err != nil {
			return err
		}
		return tgBot.Admin.Handler.(*AdminHandler).HandlerResponse.Cmd[tgBot.Text.CommonCommand.Start].Action(u)

	} else {
		return tgBot.API.SendText(u, "Invalid key! :(")
	}
}

type UnknownUnknownCmd struct{ CommandResponser }

func (r *UnknownUnknownCmd) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "Обратитесь к администратору для регистрации. Вы не авторезированный пользователь.")
	msg.ReplyMarkup = tg.NewRemoveKeyboard(false)
	_, err := tgBot.API.Send(msg)
	return err
}

/////////////////////////////////////////

type TenantStart struct {
	CommandResponser
	But keyboard.Keyboard
}

func (r *TenantStart) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Start, u.FromChat().ID))
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
	idAdmin, err := tgBot.DB.Room.GetAdmin(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	username, err := tgBot.DB.Admin.GetRepairer(idAdmin)
	if err != nil {
		return err
	}
	return tgBot.API.SendText(u, fmt.Sprintf("По техническим вопросам пишите в личные сообщения %s", username))
}

type Hot_w2 struct {
	tg.InlineKeyboardButton
	InputResponser
}

func (r *Hot_w2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
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
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
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
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
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
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
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
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
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
	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Start, u.FromChat().ID))
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type AdminCancel struct{ CommandResponser }

func (r *AdminCancel) Action(u *tg.Update) error {
	tgBot.State.Del(cache.KeyT(u.FromChat().ID))
	return tgBot.API.SendText(u, "Операция отменена")
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
}

func (r *Rooms1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "🏠Список квартирантов⋮ ")

	numbers, err := getRooms(u.FromChat().ID)
	if err != nil {
		return err
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		d.Rooms = numbers
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: st.Is, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Data: cache.AdminData{Rooms: numbers}})
	}

	msg.ReplyMarkup = keyboard.MakeInKeyboard(formatNumbers(numbers, tgBot.Text.Admin.Room2))
	_, err = tgBot.API.Send(msg)
	return err
}

type Settings1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Settings1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "🔧⋯⋯⋯⋯⋯⇐Settings⇒⋯⋯⋯⋯⋯🔧")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Room2 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *Room2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Room2) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf("💠 № 〈%s〉 :", num))
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Edit2 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *Edit2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen"); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Edit2) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "➕        OR        ✖")
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type Reminder2 struct{ InbuttonResponser }

func (r *Reminder2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Sending message."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Reminder2) Action(u *tg.Update) error {
	rooms, err := tgBot.DB.Room.Read(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}

	for _, room := range rooms {
		msg := tg.NewMessage(int64(room.IdTgTenant), "❗❗❗ Напиминаю о своевременной оплате.")
		if _, err := tgBot.API.Send(msg); err != nil {
			return err
		}
	}

	return tgBot.API.SendText(u, "✅Напоминания отправлены успешно.")
}

type Contacts2 struct{ InputResponser }

func (r *Contacts2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите @username"); err != nil {
		return err
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Admin.Settings.Contacts2, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: tgBot.Text.Admin.Settings.Contacts2})
	}

	return nil
}

type AddRoom3 struct{ InputResponser }

func (r *AddRoom3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
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
	if err := tgBot.API.AnsCallback(u, "Removing Room"); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *RemoveRoom3) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "❌Выберите какую комнату удалить➩\n🏠Список квартирантов⋮ ")

	numbers, err := getRooms(u.FromChat().ID)
	if err != nil {
		return err
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		d.RoomsDel = numbers
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: st.Is, Data: d})
	} else {
		tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Data: cache.AdminData{Rooms: numbers}})
	}
	names, data := formatNumbers(numbers, tgBot.Text.Admin.Settings.Edit.Removing4)
	names, data = append(names, []string{tgBot.Text.CommonCommand.BackBut}), append(data, []string{RemoveRoom3BackBut})
	msg.ReplyMarkup = keyboard.MakeInKeyboard(names, data)
	_, err = tgBot.API.Send(msg)
	return err
}

type Removing4 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *Removing4) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Removing4) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.NumberDel
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf("❌ Удалить комнату № 〈%s〉 и всю информацию о ней в базе?:", num))
	msg.ReplyMarkup = r.But
	_, err := tgBot.API.Send(msg)
	return err
}

type ConfirmRemove5 struct{ InbuttonResponser }

func (r *ConfirmRemove5) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Confirm?"); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ConfirmRemove5) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.NumberDel
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf("❌Комната № 〈%s〉 удалена!", num))

	if err := tgBot.DB.Scorer.Delete(database.Number(num)); err != nil {
		return err
	}
	if err := tgBot.DB.Payment.Delete(database.Number(num)); err != nil {
		return err
	}
	if err := tgBot.DB.Room.Delete(database.Number(num)); err != nil {
		return err
	}

	_, err := tgBot.API.Send(msg)
	return err
}

type ShowScorer33 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ShowScorer33) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorer33) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}
	flag := len(scorers) > MAX_SHOW_SCORER
	if flag {
		scorers = scorers[:MAX_SHOW_SCORER]
	}

	msg := tg.NewMessage(u.FromChat().ID, "🗝 **〈"+num+"〉**  ♨/💧\n"+getScorerTable(&scorers, flag))
	msg.ParseMode = tg.ModeMarkdown
	if flag {
		msg.ReplyMarkup = r.But
	}
	_, err = tgBot.API.Send(msg)
	return err
}

type ShowScorerN4 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ShowScorerN4) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing all..."); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorerN4) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}

	Emsg := tg.NewEditMessageTextAndMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, "🗝 **〈"+num+"〉**  ♨/💧\n"+getScorerTable(&scorers, false), tg.InlineKeyboardMarkup(r.But))
	Emsg.ParseMode = tg.ModeMarkdown
	_, err = tgBot.API.Send(Emsg)
	return err
}

type ShowScorerB3 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ShowScorerB3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing..."); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowScorerB3) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}
	scorers = scorers[:MAX_SHOW_SCORER]

	Emsg := tg.NewEditMessageTextAndMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, "🗝 **〈"+num+"〉**  ♨/💧\n"+getScorerTable(&scorers, true), tg.InlineKeyboardMarkup(r.But))
	Emsg.ParseMode = tg.ModeMarkdown
	_, err = tgBot.API.Send(Emsg)
	return err
}

type ShowPayment33 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ShowPayment33) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowPayment33) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	payments, err := tgBot.DB.Payment.Read(database.Number(num))
	if err != nil {
		return err
	}

	msg := tg.NewMessage(u.FromChat().ID, "🗝 **〈"+num+"〉**  🧾\n"+getPaymentTable(&payments))
	msg.ParseMode = tg.ModeMarkdown
	_, err = tgBot.API.Send(msg)
	if err != nil {
		return err
	}

	for _, payment := range payments {
		phConfig := tg.NewPhoto(u.FromChat().ID, tg.FileBytes{Name: string(payment.PayMoment), Bytes: payment.Photo})
		phConfig.Caption = string(payment.PayMoment)
		_, err = tgBot.API.Send(phConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

type ShowTenants3 struct{ InbuttonResponser }

func (r *ShowTenants3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing tenants info..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowTenants3) Action(u *tg.Update) error {
	num := ""

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	rooms, err := tgBot.DB.Room.ReadTenants(database.Number(num))
	if err != nil {
		return err
	}

	usernames := "🗝 **〈" + num + "〉**  📲\n\n"

	for _, room := range rooms {
		member, err := tgBot.API.GetChatMember(tg.GetChatMemberConfig{ChatConfigWithUser: tg.ChatConfigWithUser{ChatID: int64(room.IdTgTenant), UserID: int64(room.IdTgTenant)}})
		if err != nil {
			return nil
		}
		if member.User.UserName != "" {
			usernames += "@" + member.User.UserName + " " + member.User.FirstName + "\n"
		} else {
			usernames += member.User.FirstName + "\n"
		}
	}
	msg := tg.NewMessage(u.FromChat().ID, usernames)
	msg.ParseMode = tg.ModeMarkdown
	_, err = tgBot.API.Send(msg)
	return err
}

func getScorerTable(scorers *[]database.Scorer, flag bool) string {
	scoreTable := &strings.Builder{}
	table := tablewriter.NewWriter(scoreTable)
	//table.SetColWidth(20)
	//table.SetColMinWidth(2, 40)
	//table.SetTablePadding("*")
	//table.SetBorder(false)
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetRowSeparator("━")
	table.SetCenterSeparator("╋")
	table.SetColumnSeparator("┃") //https://unicode-table.com/ru/blocks/box-drawing/
	table.SetHeader([]string{"hot water", "cold water", "date"})
	//table.SetCaption(true, num)
	for _, score := range *scorers {
		row := []string{strconv.FormatFloat(float64(score.Hot_w)/100, 'f', -1, 64), strconv.FormatFloat(float64(score.Cold_w)/100, 'f', -1, 64), string(score.Date)}
		table.Append(row)
	}
	if flag {
		table.Append([]string{"...", "...", "..."})
	}
	table.Render()
	return "```\n" + scoreTable.String() + "```"
}

func getPaymentTable(payments *[]database.Payment) string {
	paymentTable := &strings.Builder{}
	table := tablewriter.NewWriter(paymentTable)
	table.SetRowSeparator("━")
	table.SetCenterSeparator("╋")
	table.SetColumnSeparator("┃")
	table.SetHeader([]string{"date", "amount", "pay moment"})
	for _, payment := range *payments {
		row := []string{string(payment.Date), strconv.FormatUint(uint64(payment.Amount), 10), string(payment.PayMoment)}
		table.Append(row)
	}
	table.Render()
	return "```\n" + paymentTable.String() + "```"
}
