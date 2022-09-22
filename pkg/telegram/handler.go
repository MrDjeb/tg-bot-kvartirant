package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MAX_SHOW_SCORER = 3
)

const (
	Edit2BackButInt int = iota - 256
	Room2BackButInt
	RemoveRoom3BackButInt
	Removing4BackButInt
)

var (
	Edit2BackBut       = fmt.Sprint(Edit2BackButInt)
	Room2BackBut       = fmt.Sprint(Room2BackButInt)
	RemoveRoom3BackBut = fmt.Sprint(RemoveRoom3BackButInt)
	Removing4BackBut   = fmt.Sprint(Removing4BackButInt)
)

type Handler interface {
	Callback(u *tg.Update) error
	Command(u *tg.Update) error
	Photo(u *tg.Update) error
	Message(u *tg.Update) error
	New()
}

type HandlerResponse struct {
	Cmd map[string]CommandResponser
	Mes map[string]MessageResponser
	But map[string]ButtonResponser
	Inb map[string]InbuttonResponser
	Inp map[string]InputResponser
	Bck map[string]BackResponser
}

type CommandResponser interface {
	Action(u *tg.Update) error
}

type MessageResponser interface {
	Action(u *tg.Update) error
}

type ButtonResponser interface {
	Action(u *tg.Update) error
}

type InputResponser interface {
	Callback(u *tg.Update) error
	HandleInput(u *tg.Update) error
}

type InbuttonResponser interface {
	Callback(u *tg.Update) error
	Action(u *tg.Update) error
}

type BackResponser func(u *tg.Update) error

///////////////////////////////////////////////////

type UnknownHandler struct{ HandlerResponse }

func (h *UnknownHandler) New() {
	h.Cmd = map[string]CommandResponser{
		tgBot.Text.CommonCommand.Start:   &UnknownStart{},
		tgBot.Text.CommonCommand.GodMode: &GodMode{},
		tgBot.Text.CommonCommand.Unknown: &UnknownUnknownCmd{},
	}
}

func (h *UnknownHandler) Callback(u *tg.Update) error {
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *UnknownHandler) Command(u *tg.Update) error {
	if cmd, ok := h.Cmd[u.Message.Command()]; ok {
		return cmd.Action(u)
	}
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *UnknownHandler) Photo(u *tg.Update) error {
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *UnknownHandler) Message(u *tg.Update) error {
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

// /////////////////
type TenantHandler struct{ HandlerResponse }

func (h *TenantHandler) New() {
	TC := tgBot.Text.CommonCommand
	TT := tgBot.Text.Tenant
	h.Cmd = map[string]CommandResponser{
		TC.Start:   &TenantStart{But: keyboard.MakeKeyboard([]string{TT.Water1, TT.Receipt1}, []string{TT.Report1})},
		TC.Cancel:  &TenantCancel{},
		TC.Unknown: &TenantUnknownCmd{},
	}
	h.Mes = map[string]MessageResponser{
		tgBot.Text.CommonMessage.Hi: &TenantHi{},
		TC.Unknown:                  &TenantUnknownMes{},
	}
	h.But = map[string]ButtonResponser{
		TT.Water1:   &Water1{But: keyboard.MakeInKeyboard([][]string{{TT.Water.Hot_w2, TT.Water.Cold_w2}, {TT.Water.Choose_month}}, [][]string{{TT.Water.Hot_w2, TT.Water.Cold_w2}, {TT.Water.Choose_month}})},
		TT.Receipt1: &Receipt1{But: keyboard.MakeInKeyboard([][]string{{TT.Receipt.Month2, TT.Receipt.Amount2, TT.Receipt.Receipt2}}, [][]string{{TT.Receipt.Month2, TT.Receipt.Amount2, TT.Receipt.Receipt2}})},
		TT.Report1:  &Report1{},
	}
	h.Inb = map[string]InbuttonResponser{
		TT.Water.Choose_month:   &ChooseMonth{Prefix: TT.Water.Month_prefix},
		TT.Water.Month_prefix:   &WaterM1{But: keyboard.MakeInKeyboard([][]string{{TT.Water.Hot_w2, TT.Water.Cold_w2}, {TT.Water.Choose_month}}, [][]string{{TT.Water.Hot_w2, TT.Water.Cold_w2}, {TT.Water.Choose_month}})},
		TT.Receipt.Month2:       &ChooseMonth{Prefix: TT.Receipt.Month_prefix},
		TT.Receipt.Month_prefix: &Month2{},
	}
	h.Inp = map[string]InputResponser{
		TT.Water.Cold_w2:    &Cold_w2{},
		TT.Water.Hot_w2:     &Hot_w2{},
		TT.Receipt.Amount2:  &Amount2{},
		TT.Receipt.Receipt2: &Receipt2{},
	}
}

func (h *TenantHandler) Callback(u *tg.Update) error {
	if inp, ok := h.Inp[u.CallbackQuery.Data]; ok {
		return inp.Callback(u)
	}
	if Inb, ok := h.Inb[u.CallbackQuery.Data]; ok {
		return Inb.Callback(u)
	}

	if strings.Contains(u.CallbackQuery.Data, keyboard.DEL) {
		prefix, suffix := strings.Split(u.CallbackQuery.Data, keyboard.DEL)[0], strings.Split(u.CallbackQuery.Data, keyboard.DEL)[1]
		if inb, ok := h.Inb[prefix]; ok {
			month, err := strconv.Atoi(strings.Split(suffix, "-")[1])
			if err != nil {
				return err
			}
			switch prefix {
			case tgBot.Text.Tenant.Water.Month_prefix:
				if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
					d.ScoreDate = uint8(month)
					tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
				} else {
					tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{ScoreDate: uint8(month)})
				}
			case tgBot.Text.Tenant.Receipt.Month_prefix:
				if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
					d.PaymentDate = uint8(month)
					tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
				} else {
					tgBot.Tenant.Cache.Put(u.FromChat().ID, TenantData{PaymentDate: uint8(month)})
				}
			}
			return inb.Callback(u)
		}
	}

	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *TenantHandler) Command(u *tg.Update) error {
	if cmd, ok := h.Cmd[u.Message.Command()]; ok {
		return cmd.Action(u)
	}
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *TenantHandler) Photo(u *tg.Update) error {
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && d.Is == tgBot.Text.Buttons.Tenant.Receipt.Receipt2 {
		return h.Inp[tgBot.Text.Tenant.Receipt.Receipt2].HandleInput(u)
	} else if ok {
		return tgBot.API.SendText(u, "Сейчас мне не нужно фото.")
	} else {
		return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
	}
}

func (h *TenantHandler) Message(u *tg.Update) error {
	if mes, ok := h.Mes[u.Message.Text]; ok {
		return mes.Action(u)
	}
	if but, ok := h.But[u.Message.Text]; ok {
		return but.Action(u)
	}

	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
		if d.Is == tgBot.Text.Tenant.Receipt.Receipt2 {
			return tgBot.API.SendText(u, "Пришлите фото.")
		}
		if inp, ok := h.Inp[d.Is]; ok {
			return inp.HandleInput(u)
		}
	}
	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

// /////////////////////
type AdminHandler struct{ HandlerResponse }

func (h *AdminHandler) New() {
	TC := tgBot.Text.CommonCommand
	TA := tgBot.Text.Admin
	h.Cmd = map[string]CommandResponser{
		TC.Start:   &AdminStart{But: keyboard.MakeKeyboard([]string{TA.Rooms1, TA.Settings1})},
		TC.Cancel:  &AdminCancel{},
		TC.Unknown: &AdminUnknownCmd{},
	}
	h.Mes = map[string]MessageResponser{
		tgBot.Text.CommonMessage.Hi: &AdminHi{},
		TC.Unknown:                  &AdminUnknownMes{},
	}
	h.But = map[string]ButtonResponser{
		TA.Rooms1:    &Rooms1{},
		TA.Settings1: &Settings1{But: keyboard.MakeInKeyboard([][]string{{TA.Settings.Edit2}, {TA.Settings.Contacts2}, {TA.Settings.Reminder2}}, [][]string{{TA.Settings.Edit2}, {TA.Settings.Contacts2}, {TA.Settings.Reminder2}})},
	}
	h.Inb = map[string]InbuttonResponser{
		TA.Room2:                                 &Room2{But: keyboard.MakeInKeyboard([][]string{{TA.Room.ShowScorer33, TA.Room.ShowPayment33}, {TA.Room.ShowTenants3}, {TC.BackBut}}, [][]string{{TA.Room.ShowScorer33, TA.Room.ShowPayment33}, {TA.Room.ShowTenants3}, {Room2BackBut}})},
		TA.Room.ShowScorer33:                     &ShowScorer33{But: keyboard.MakeInKeyboard([][]string{{TA.Room.ShowScorerN4}}, [][]string{{TA.Room.ShowScorerN4}})},
		TA.Room.ShowScorerN4:                     &ShowScorerN4{But: keyboard.MakeInKeyboard([][]string{{TA.Room.ShowScorerB3}}, [][]string{{TA.Room.ShowScorerB3}})},
		TA.Room.ShowScorerB3:                     &ShowScorerB3{But: keyboard.MakeInKeyboard([][]string{{TA.Room.ShowScorerN4}}, [][]string{{TA.Room.ShowScorerN4}})},
		TA.Room.ShowPayment33:                    &ShowPayment33{},
		TA.Room.ShowTenants3:                     &ShowTenants3{},
		TA.Settings.Edit2:                        &Edit2{But: keyboard.MakeInKeyboard([][]string{{TA.Settings.Edit.AddRoom3, TA.Settings.Edit.RemoveRoom3}, {TC.BackBut}}, [][]string{{TA.Settings.Edit.AddRoom3, TA.Settings.Edit.RemoveRoom3}, {Edit2BackBut}})},
		TA.Settings.Reminder2:                    &Reminder2{},
		TA.Settings.Edit.RemoveRoom3:             &RemoveRoom3{},
		TA.Settings.Edit.Removing4:               &Removing4{But: keyboard.MakeInKeyboard([][]string{{TA.Settings.Edit.Removing.ConfirmRemove5}, {TC.BackBut}}, [][]string{{TA.Settings.Edit.Removing.ConfirmRemove5}, {Removing4BackBut}})},
		TA.Settings.Edit.Removing.ConfirmRemove5: &ConfirmRemove5{},
	}
	h.Inp = map[string]InputResponser{
		TA.Settings.Contacts2:     &Contacts2{},
		TA.Settings.Edit.AddRoom3: &AddRoom3{},
	}
	h.Bck = map[string]BackResponser{
		Edit2BackBut:       NewBackResponser(TA.Settings1),
		Room2BackBut:       NewBackResponser(TA.Rooms1),
		RemoveRoom3BackBut: NewBackResponser(TA.Settings.Edit2),
		Removing4BackBut:   NewBackResponser(TA.Settings.Edit.RemoveRoom3),
	}
}

func (h *AdminHandler) Callback(u *tg.Update) error {
	if bck, ok := h.Bck[u.CallbackQuery.Data]; ok {
		return bck(u)
	}
	if inp, ok := h.Inp[u.CallbackQuery.Data]; ok {
		return inp.Callback(u)
	}
	if inb, ok := h.Inb[u.CallbackQuery.Data]; ok {
		return inb.Callback(u)
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		num := u.CallbackQuery.Data[strings.Index(u.CallbackQuery.Data, keyboard.DEL)+1:]
		if isRoom(num, d.Rooms) {
			prefix := u.CallbackQuery.Data[:strings.Index(u.CallbackQuery.Data, keyboard.DEL)]

			if inb, ok := h.Inb[prefix]; ok {
				switch prefix {
				case tgBot.Text.Admin.Room2:
					d.Number = num
				case tgBot.Text.Admin.Settings.Edit.Removing4:
					d.NumberDel = num
				}
				tgBot.Admin.Cache.Put(u.FromChat().ID, d)
				return inb.Callback(u)
			}
		}
	}

	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *AdminHandler) Command(u *tg.Update) error {
	if cmd, ok := h.Cmd[u.Message.Command()]; ok {
		return cmd.Action(u)
	}
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *AdminHandler) Photo(u *tg.Update) error {
	return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
}

func (h *AdminHandler) Message(u *tg.Update) error {
	if but, ok := h.But[u.Message.Text]; ok {
		return but.Action(u)
	}
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok && d.Is != "" {
		if inp, ok := h.Inp[d.Is]; ok {
			return inp.HandleInput(u)
		}
	}
	if mes, ok := h.Mes[u.Message.Text]; ok {
		return mes.Action(u)
	}
	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func getRooms(idTg int64) ([]string, error) {
	rooms, err := tgBot.DB.Room.Read(database.TelegramID(idTg))
	if err != nil {
		return nil, err
	}
	var numbers []string
	for _, room := range rooms {
		numbers = append(numbers, string(room.Number))
	}
	return numbers, nil
}

func isRoom(number string, numbers []string) bool {
	for _, num := range numbers {
		if num == number {
			return true
		}
	}
	return false
}

func NewBackResponser(t string) BackResponser {
	return func(u *tg.Update) error {
		h := tgBot.Admin.Handler.(*AdminHandler).HandlerResponse
		if err := tgBot.API.AnsCallback(u, "BACK"); err != nil {
			return err
		}
		if err := tgBot.API.DelMes(u); err != nil {
			return err
		}

		but, ok := h.But[t]
		if ok {
			return but.Action(u)
		}
		inb, ok := h.Inb[t]
		if ok {
			return inb.Action(u)
		}
		inp, ok := h.Inp[t]
		if ok {
			return inp.HandleInput(u)
		}
		mes, ok := h.Mes[t]
		if ok {
			return mes.Action(u)
		}
		cmd, ok := h.Cmd[t]
		if ok {
			return cmd.Action(u)
		}
		return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
	}
}
