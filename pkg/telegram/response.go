package telegram

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_start+"(err1)")
	}
	token := string(byteToken)
	idAdmin, err := strconv.ParseInt(token[32:], 10, 64)
	if err != nil { //error broke
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_start+"(err2)")
	}

	d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(idAdmin)
	if !ok {
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_start+"(err3),Кэш с таким idAdmin пуст.")
	}
	number, ok := d.AddingRooms[token]
	if !ok {
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_start+"(err4)")
	}
	delete(d.AddingRooms, token)
	tgBot.Admin.Cache.Put(idAdmin, d)

	room := &database.Room{
		IdTgAdmin:  database.TelegramID(idAdmin),
		IdTgTenant: database.TelegramID(u.FromChat().ID),
		Number:     database.Number(number),
	}
	if err := tgBot.DB.Room.Insert(*room); err != nil {
		return err
	}

	msg := tg.NewMessage(idAdmin, fmt.Sprintf("🔗 %s %s успешно привязался(-aсь) к комнате № 〈%s〉", u.SentFrom().FirstName, u.SentFrom().UserName, number))
	_, err = tgBot.API.Send(msg)
	if err != nil {
		return err
	}

	if err := tgBot.DB.Tenant.Insert(database.Tenant{IdTg: database.TelegramID(u.FromChat().ID)}); err != nil {
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
	msg := tg.NewMessage(u.FromChat().ID, "Обратитесь к администратору для регистрации. Вы не авторизированный пользователь.")
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
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		d.Erase()
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
	}
	msg := tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Cancel)
	msg.ParseMode = tg.ModeMarkdownV2
	_, err := tgBot.API.Send(msg)
	return err
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
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		if d.Payment[0] || d.Payment[1] || d.Payment[2] || d.Score[0] || d.Score[1] {
			return tgBot.API.SendText(u, "Продолжите внесение данных.")
		}
	}
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
}

type Water1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Water1) Action(u *tg.Update) error {
	var msg tg.MessageConfig

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && (d.Score[0] != d.Score[1]) { //xor

		msg = tg.NewMessage(u.Message.Chat.ID, tgBot.Text.Response.Water1_sec)
		var inlineButtons []tg.InlineKeyboardButton
		for i := range d.Score {
			if !d.Score[i] {
				inlineButtons = append(inlineButtons, r.But.InlineKeyboard[0][i])
			}
		}
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(inlineButtons)
	} else {
		date := "текущий"
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && d.ScoreDate != 0 {
			date = getAverageDate(d.ScoreDate)
		}
		msg = tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Water1_first, date))
		msg.ReplyMarkup = r.But
	}

	_, err := tgBot.API.Send(msg)
	return err
}

type Receipt1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Receipt1) Action(u *tg.Update) error {
	var msg tg.MessageConfig

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		switch tint(d.Payment[0]) + tint(d.Payment[1]) + tint(d.Payment[2]) {
		case 0:
			msg = tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Receipt1_first)
		case 1:
			msg = tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Receipt1_sec)
		case 2:
			msg = tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Receipt1_third)
		}
		var inlineButtons []tg.InlineKeyboardButton
		for i := range d.Payment {
			if !d.Payment[i] {
				inlineButtons = append(inlineButtons, r.But.InlineKeyboard[0][i])
			}
		}
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(inlineButtons)
	} else {
		msg = tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Receipt1_first)
		msg.ReplyMarkup = r.But
	}

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

	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Report1_info, username, username))
	msg.ParseMode = tg.ModeMarkdown
	if _, err = tgBot.API.Send(msg); err != nil {
		return err
	}

	msg = tg.NewMessage(u.FromChat().ID, fmt.Sprintf("@%s", username))
	msg.Entities = append(msg.Entities, tg.MessageEntity{Type: "mention"})
	if _, err = tgBot.API.Send(msg); err != nil {
		return err
	}
	return nil
}

type Hot_w2 struct {
	InputResponser
}

func (r *Hot_w2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, tgBot.Text.Response.Water2_inp); err != nil {
		return err
	}

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Tenant.Water.Hot_w2
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{Is: tgBot.Text.Tenant.Water.Hot_w2})
	}

	return nil
}

type Cold_w2 struct {
	InputResponser
}

func (r *Cold_w2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, tgBot.Text.Response.Water2_inp); err != nil {
		return err
	}

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Tenant.Water.Cold_w2
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{Is: tgBot.Text.Tenant.Water.Cold_w2})
	}
	return nil
}

type Amount2 struct {
	InputResponser
}

func (r *Amount2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, tgBot.Text.Response.Amount2_inp); err != nil {
		return err
	}

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Tenant.Receipt.Amount2
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{Is: tgBot.Text.Tenant.Receipt.Amount2})
	}
	return nil
}

type Receipt2 struct {
	InputResponser
}

func (r *Receipt2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, "Прикрепите изображение квитанции, подтверждающее факт оплаты."); err != nil {
		return err
	}

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Tenant.Receipt.Receipt2
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{Is: tgBot.Text.Tenant.Receipt.Receipt2})
	}
	return nil
}

type Month2 struct {
	InbuttonResponser
	Prefix string
}

func (r *Month2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Selected..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)

}

func (r *Month2) Action(u *tg.Update) error {
	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}

	d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
	}
	d.Payment[0] = true

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentDate)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Receipt2_saved, payment.Amount, payment.Date)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		if err := tgBot.API.SendText(u, "Месяц успешно добавлен."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Receipt1].Action(u); err != nil {
			return err
		}
	}

	return nil
}

type ChooseMonth struct {
	InbuttonResponser
	Prefix string
}

func (r *ChooseMonth) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ChooseMonth) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, "🌙 Выберите месяц:")

	/*InName, InData := keyboard.MakeFormatMonth(r.Prefix)
	for i := range InName {
		for j := range InName[i] {
				mon, _ := strconv.Atoi(InName[i][j])
				fmt.Println(getAverageDate(uint8(mon))[:4])
				yer, _ := strconv.Atoi(getAverageDate(uint8(mon))[:4])
				if time.Now().Year() < yer {
					InName[i][j] += "+"
				} else if time.Now().Year() > yer {
					InName[i][j] += "-"
				}
			}
		}*/

	msg.ReplyMarkup = keyboard.MakeInKeyboard(getFormatCalendar(r.Prefix))
	_, err := tgBot.API.Send(msg)
	return err
}

type WaterM1 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *WaterM1) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *WaterM1) Action(u *tg.Update) error {
	date := "текущий"
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && d.ScoreDate != 0 {
		date = getAverageDate(d.ScoreDate)
	}
	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Water1_first, date))
	msg.ReplyMarkup = r.But

	_, err := tgBot.API.Send(msg)
	return err
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
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = ""
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	}
	msg := tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Cancel)
	msg.ParseMode = tg.ModeMarkdownV2
	_, err := tgBot.API.Send(msg)
	return err
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
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)

}

type Rooms1 struct {
	ButtonResponser
}

func (r *Rooms1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Rooms1_list)

	fNum, fData, err := getFormatRooms(u.FromChat().ID, tgBot.Text.Admin.Room2)
	if err != nil {
		return err
	}

	if len(fData) == 0 {
		return tgBot.API.SendText(u, tgBot.Text.Response.Rooms1_nil)
	}

	msg.ReplyMarkup = keyboard.MakeInKeyboard(fNum, fData)
	_, err = tgBot.API.Send(msg)
	return err
}

type Settings1 struct {
	ButtonResponser
	But keyboard.InKeyboard
}

func (r *Settings1) Action(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, tgBot.Text.Response.Settings1)
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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf(tgBot.Text.Response.Room2, num))
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

type Contacts2 struct{ InputResponser }

func (r *Contacts2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Start inputing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, "Введите @username"); err != nil {
		return err
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Admin.Settings.Contacts2
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Admin.Cache.Put(u.FromChat().ID, AdminData{Is: tgBot.Text.Admin.Settings.Contacts2})
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
	if err := tgBot.API.SendInputText(u, "Введите номер комнаты, который хотите добавить"); err != nil {
		return err
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Admin.Settings.Edit.AddRoom3
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Admin.Cache.Put(u.FromChat().ID, AdminData{Is: tgBot.Text.Admin.Settings.Edit.AddRoom3,
			AddingRooms: make(map[string]string)})
	}
	return nil
}

type Reminder2 struct {
	InputResponser
	But keyboard.InKeyboard
}

func (r *Reminder2) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Sending message."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, "Введите сообщение, копию которого хотите разослать:"); err != nil {
		return err
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Admin.Settings.Reminder2
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Admin.Cache.Put(u.FromChat().ID, AdminData{Is: tgBot.Text.Admin.Settings.Reminder2})
	}
	return nil
}

type ReminderEdit3 struct {
	InputResponser
	But keyboard.InKeyboard
}

func (r *ReminderEdit3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Sending message."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	if err := tgBot.API.SendInputText(u, "Введите сообщение, копию которого хотите разослать:"); err != nil {
		return err
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Admin.Settings.Reminder2
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Admin.Cache.Put(u.FromChat().ID, AdminData{Is: tgBot.Text.Admin.Settings.Reminder2})
	}
	return nil
}

type ReminderSend3 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *ReminderSend3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Choosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ReminderSend3) Action(u *tg.Update) error {
	text := ""
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		text = d.RemindText
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	rooms, err := tgBot.DB.Room.Read(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	for _, room := range rooms {
		msg := tg.NewMessage(int64(room.IdTgTenant), text)
		if _, err := tgBot.API.GetChat(tg.ChatInfoConfig{ChatConfig: tg.ChatConfig{ChatID: int64(room.IdTgTenant)}}); err == nil {
			if _, err := tgBot.API.Send(msg); err != nil {
				return err
			}
		}
	}

	return tgBot.API.SendText(u, "✅Напоминания успешно отправлены.")
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
	msg := tg.NewMessage(u.FromChat().ID, "❌Выберите, какую комнату хотите удалить➩\n")

	fNum, fData, err := getFormatRooms(u.FromChat().ID, tgBot.Text.Admin.Settings.Edit.Removing4)
	if err != nil {
		return err
	}

	if len(fData) == 0 {
		return tgBot.API.SendText(u, tgBot.Text.Response.Rooms1_nil)
	}

	fNum, fData = append(fNum, []string{tgBot.Text.CommonCommand.BackBut}), append(fData, []string{RemoveRoom3BackBut})
	msg.ReplyMarkup = keyboard.MakeInKeyboard(fNum, fData)
	_, err = tgBot.API.Send(msg)
	return err
}

type Removing4 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

func (r *Removing4) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Chosen..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *Removing4) Action(u *tg.Update) error {
	num := ""

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.NumberDel
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}
	if err := DeleteRoom(database.Number(num)); err != nil {
		return err
	}
	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf("❌Комната № 〈%s〉 удалена!", num))

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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}
	if len(scorers) == 0 {
		return tgBot.API.SendText(u, "Пользователь пока ничего не вносил.")
	}
	sort.Slice(scorers, func(i, j int) bool { return scorers[i].Date > scorers[j].Date })

	flag := len(scorers) > tgBot.Text.Constants.MAX_SHOW_SCORER
	if flag {
		scorers = scorers[:tgBot.Text.Constants.MAX_SHOW_SCORER]
	}

	msg := tg.NewMessage(u.FromChat().ID, "🗝 **〈"+num+"〉**  ♨/💧\n"+getScorerTable(&scorers, flag))
	msg.ParseMode = tg.ModeMarkdown
	//fmt.Println("MessagwEntity: ", msg.Entities)
	//msg.Entities = append(msg.Entities, tg.MessageEntity{Type: "code"})
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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}
	sort.Slice(scorers, func(i, j int) bool { return scorers[i].Date > scorers[j].Date })

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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	scorers, err := tgBot.DB.Scorer.Read(database.Number(num))
	if err != nil {
		return err
	}
	sort.Slice(scorers, func(i, j int) bool { return scorers[i].Date > scorers[j].Date })
	scorers = scorers[:tgBot.Text.Constants.MAX_SHOW_SCORER]

	Emsg := tg.NewEditMessageTextAndMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, "🗝 **〈"+num+"〉**  ♨/💧\n"+getScorerTable(&scorers, true), tg.InlineKeyboardMarkup(r.But))
	Emsg.ParseMode = tg.ModeMarkdown
	_, err = tgBot.API.Send(Emsg)
	return err
}

type ShowPayment33 struct {
	InbuttonResponser
	But    keyboard.InKeyboard
	Prefix string
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
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	payments, err := tgBot.DB.Payment.Read(database.Number(num))
	if err != nil {
		return err
	}
	if len(payments) == 0 {
		return tgBot.API.SendText(u, "Пользователь пока ничего не вносил.")
	}

	fNum, fData := getFormatPayment(r.Prefix, &payments)
	flag := len(fNum) > tgBot.Text.Constants.MAX_SHOW_PAYMENT
	if flag {
		fNum = fNum[:tgBot.Text.Constants.MAX_SHOW_PAYMENT]
		fData = fData[:tgBot.Text.Constants.MAX_SHOW_PAYMENT]
	}

	buttons := keyboard.MakeInKeyboard(fNum, fData)
	if flag {
		buttons.InlineKeyboard = append(buttons.InlineKeyboard, r.But.InlineKeyboard[0])
	}
	msg := tg.NewMessage(u.FromChat().ID, "🗝 **〈"+num+"〉**  🧾\n")
	msg.ReplyMarkup = buttons
	msg.ParseMode = tg.ModeMarkdownV2
	_, err = tgBot.API.Send(msg)
	return err
}

type ShowPaymentN4 struct {
	InbuttonResponser
	But    keyboard.InKeyboard
	Prefix string
}

func (r *ShowPaymentN4) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing all..."); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowPaymentN4) Action(u *tg.Update) error {
	num := ""
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	payments, err := tgBot.DB.Payment.Read(database.Number(num))
	if err != nil {
		return err
	}

	buttons := keyboard.MakeInKeyboard(getFormatPayment(r.Prefix, &payments))
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, r.But.InlineKeyboard[0])

	Emsg := tg.NewEditMessageTextAndMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, "🗝 **〈"+num+"〉**  🧾\n", tg.InlineKeyboardMarkup(buttons))
	Emsg.ParseMode = tg.ModeMarkdownV2
	_, err = tgBot.API.Send(Emsg)
	return err
}

type ShowPaymentB3 struct {
	InbuttonResponser
	But    keyboard.InKeyboard
	Prefix string
}

func (r *ShowPaymentB3) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing..."); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowPaymentB3) Action(u *tg.Update) error {
	num := ""
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}
	payments, err := tgBot.DB.Payment.Read(database.Number(num))
	if err != nil {
		return err
	}

	fNum, fData := getFormatPayment(r.Prefix, &payments)
	flag := len(fNum) > tgBot.Text.Constants.MAX_SHOW_PAYMENT
	if flag {
		fNum = fNum[:tgBot.Text.Constants.MAX_SHOW_PAYMENT]
		fData = fData[:tgBot.Text.Constants.MAX_SHOW_PAYMENT]
	}

	buttons := keyboard.MakeInKeyboard(fNum, fData)
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, r.But.InlineKeyboard[0])

	Emsg := tg.NewEditMessageTextAndMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, "🗝 **〈"+num+"〉**  🧾\n", tg.InlineKeyboardMarkup(buttons))
	Emsg.ParseMode = tg.ModeMarkdownV2
	_, err = tgBot.API.Send(Emsg)
	return err
}

type ShowPayment struct {
	InbuttonResponser
}

func (r *ShowPayment) Callback(u *tg.Update) error {
	if err := tgBot.API.AnsCallback(u, "Showing..."); err != nil {
		return err
	}
	if err := tgBot.API.DelMes(u); err != nil {
		return err
	}
	return r.Action(u)
}

func (r *ShowPayment) Action(u *tg.Update) error {
	date, num := "", ""
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		date, num = d.ShowPayment, d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	payments, err := tgBot.DB.Payment.Read(database.Number(num))
	if err != nil {
		return err
	}

	for _, payment := range payments {
		if date == string(payment.Date) {
			phConfig := tg.NewPhoto(u.FromChat().ID, tg.FileBytes{Name: string(payment.PayMoment), Bytes: payment.Photo})
			phConfig.Caption = fmt.Sprintf("Дата внесения: %s, сумма: %d₽", string(payment.PayMoment), uint(payment.Amount))

			_, err = tgBot.API.Send(phConfig)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*func (r *ShowPayment33) Action(u *tg.Update) error {
	num := ""
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
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
}*/

type ShowTenants3 struct {
	InbuttonResponser
	But keyboard.InKeyboard
}

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

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num = d.Number
	} else {
		return errors.New("gets the nil data from cache, can't do function")
	}

	rooms, err := tgBot.DB.Room.ReadRooms(database.Number(num))
	fmt.Println(rooms)
	if err != nil {
		return err
	}

	usernames := "🗝 **〈" + num + "〉**  📲\n\n⇣Список пользователей телеграм с установленным username\n\n"

	for _, room := range rooms {
		if _, err := tgBot.API.GetChat(tg.ChatInfoConfig{ChatConfig: tg.ChatConfig{ChatID: int64(room.IdTgTenant)}}); err == nil {
			member, err := tgBot.API.GetChatMember(tg.GetChatMemberConfig{ChatConfigWithUser: tg.ChatConfigWithUser{ChatID: int64(room.IdTgTenant), UserID: int64(room.IdTgTenant)}})
			if err != nil {
				return nil
			}
			usernames += fmt.Sprintf("‣ [%s](tg://user?id=%d)    `%d`\n", member.User.FirstName, member.User.ID, room.IdTgTenant)
		} else {
			usernames += fmt.Sprintf("‣ (tg:ID)    `%d`\n", room.IdTgTenant)
		}
	}
	msg := tg.NewMessage(u.FromChat().ID, usernames)
	msg.ParseMode = tg.ModeMarkdown
	msg.ReplyMarkup = r.But
	_, err = tgBot.API.Send(msg)
	return err
}

type RemoveTenants4 struct{ InputResponser }

func (r *RemoveTenants4) Callback(u *tg.Update) error { /////////////////////////neeed

	if err := tgBot.API.AnsCallback(u, "Sending message."); err != nil {
		return err
	}

	Emsg := tg.NewEditMessageReplyMarkup(u.FromChat().ID, u.CallbackQuery.Message.MessageID, tg.InlineKeyboardMarkup{InlineKeyboard: make([][]tg.InlineKeyboardButton, 0)})
	if _, err := tgBot.API.Send(Emsg); err != nil {
		return err
	}

	if err := tgBot.API.SendInputText(u, "Скопируйте Telegram ID нужного квартиранта (выберите в списке и нажмите), затем вставьте и отправьте его:"); err != nil {
		return err
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = tgBot.Text.Admin.Room.RemoveTenants4
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	} else {
		tgBot.Admin.Cache.Put(u.FromChat().ID, AdminData{Is: tgBot.Text.Admin.Room.RemoveTenants4})
	}
	return nil
}

func DeleteRoom(num database.Number) error {
	rooms, err := tgBot.DB.Room.ReadRooms(database.Number(num))
	if err != nil {
		return err
	}
	for _, r := range rooms {
		if err := tgBot.DB.Tenant.Delete(r.IdTgTenant); err != nil {
			return err
		}
	}
	if err := tgBot.DB.Scorer.Delete(database.Number(num)); err != nil {
		return err
	}
	if err := tgBot.DB.Payment.Delete(database.Number(num)); err != nil {
		return err
	}
	if err := tgBot.DB.Room.Delete(database.Number(num)); err != nil {
		return err
	}
	return err
}

func DeleteTenant(tgid database.TelegramID) error {
	/*num, err := tgBot.DB.Room.GetRoom(tgid)
	if err != nil {
		return err
	}

	rooms, err := tgBot.DB.Room.ReadRooms(database.Number(num))
	if err != nil {
		return err
	}
	if len(rooms) == 1 {
		return DeleteRoom(num)
	}*/
	if err := tgBot.DB.Room.DeleteTenant(tgid); err != nil {
		return err
	}
	return tgBot.DB.Tenant.Delete(tgid)
}

func getScorerTable(scorers *[]database.Scorer, flag bool) string {
	scoreTable := &strings.Builder{}
	table := tablewriter.NewWriter(scoreTable)
	//table.SetColWidth(20)
	//table.SetColMinWidth(2, 40)
	//table.SetTablePadding("*")
	//table.SetBorder(false)
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	//table.SetRowSeparator("-")
	//table.SetCenterSeparator("+")
	//table.SetColumnSeparator("|") //https://unicode-table.com/ru/blocks/box-drawing/
	table.SetHeader([]string{"hot", "cold", "date"})
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

/*
func getPaymentTable(payments *[]database.Payment) string {
	paymentTable := &strings.Builder{}
	table := tablewriter.NewWriter(paymentTable)
	//table.SetRowSeparator("━")
	//table.SetCenterSeparator("╋")
	//table.SetColumnSeparator("┃")

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.SetHeader([]string{"date", "amount", "moment"})
	for _, payment := range *payments {
		row := []string{string(payment.Date), strconv.FormatUint(uint64(payment.Amount), 10), string(payment.PayMoment)}
		table.Append(row)
	}
	table.Render()
	return "```\n" + paymentTable.String() + "```"
}


func markImage(fotoByte []byte) tg.PhotoConfig {
	img, _, _ := image.Decode(bytes.NewReader(fotoByte))

	return tg.NewPhoto(u.FromChat().ID, tg.FileBytes{Name: string(payment.PayMoment), Bytes: payment.Photo})
}*/
