package telegram

import (
	"fmt"
	"strings"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DEL             = "$"
	MAX_SHOW_SCORER = 5
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

func NewUser(h Handler) User {
	h.New()
	return User{h}
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
	cmd, ok := h.Cmd[u.Message.Command()]
	if ok {
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

///////////////////
type TenantHandler struct{ HandlerResponse }

func (h *TenantHandler) New() {
	b := keyboard.NewButtons().Tenant
	t := tgBot.Text.Tenant
	h.Cmd = map[string]CommandResponser{
		tgBot.Text.CommonCommand.Start:   &TenantStart{But: keyboard.MakeKeyboard([]string{t.Water1, t.Receipt1}, []string{t.Report1})},
		tgBot.Text.CommonCommand.Cancel:  &TenantCancel{},
		tgBot.Text.CommonCommand.Unknown: &TenantUnknownCmd{},
	}
	h.Mes = map[string]MessageResponser{
		tgBot.Text.CommonMessage.Hi:      &TenantHi{},
		tgBot.Text.CommonCommand.Unknown: &TenantUnknownMes{},
	}
	h.But = map[string]ButtonResponser{
		tgBot.Text.Tenant.Water1:   &Water1{But: b.Water},
		tgBot.Text.Tenant.Receipt1: &Receipt1{But: b.Receipt},
		tgBot.Text.Tenant.Report1:  &Report1{},
	}
	h.Inp = map[string]InputResponser{
		tgBot.Text.Tenant.Water.Cold_w2:    &Cold_w2{},
		tgBot.Text.Tenant.Water.Hot_w2:     &Hot_w2{},
		tgBot.Text.Tenant.Receipt.Amount2:  &Amount2{},
		tgBot.Text.Tenant.Receipt.Month2:   &Month2{},
		tgBot.Text.Tenant.Receipt.Receipt2: &Receipt2{},
	}
}

func (h *TenantHandler) Callback(u *tg.Update) error {
	inp, ok := h.Inp[u.CallbackQuery.Data]
	if ok {
		return inp.Callback(u)
	}
	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *TenantHandler) Command(u *tg.Update) error {
	cmd, ok := h.Cmd[u.Message.Command()]
	if ok {
		return cmd.Action(u)
	}
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *TenantHandler) Photo(u *tg.Update) error {
	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok && st.Is == tgBot.Text.Buttons.Tenant.Receipt.Receipt2 {
		return h.Inp[tgBot.Text.Tenant.Receipt.Receipt2].HandleInput(u)
	} else if ok {
		return tgBot.API.SendText(u, "Сейчас мне не нужно фото.")
	} else {
		return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
	}
}

func (h *TenantHandler) Message(u *tg.Update) error {
	mes, ok := h.Mes[u.Message.Text]
	if ok {
		return mes.Action(u)
	}

	but, ok := h.But[u.Message.Text]
	if ok {
		return but.Action(u)
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		if st.Is == tgBot.Text.Tenant.Receipt.Receipt2 {
			return tgBot.API.SendText(u, "Пришлите фото.")
		}
		return h.Inp[st.Is].HandleInput(u)
	}
	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

///////////////////////
type AdminHandler struct{ HandlerResponse }

func (h *AdminHandler) New() {
	B := keyboard.NewButtons().Admin
	TC := tgBot.Text.CommonCommand
	TA := tgBot.Text.Admin
	h.Cmd = map[string]CommandResponser{
		TC.Start:   &AdminStart{But: B.Keyboard},
		TC.Cancel:  &AdminCancel{},
		TC.Unknown: &AdminUnknownCmd{},
	}
	h.Mes = map[string]MessageResponser{
		tgBot.Text.CommonMessage.Hi: &AdminHi{},
		TC.Unknown:                  &AdminUnknownMes{},
	}
	h.But = map[string]ButtonResponser{
		TA.Rooms1:    &Rooms1{},
		TA.Settings1: &Settings1{But: B.Settings},
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
	bck, ok := h.Bck[u.CallbackQuery.Data]
	if ok {
		return bck(u)
	}

	inp, ok := h.Inp[u.CallbackQuery.Data]
	if ok {
		return inp.Callback(u)
	}

	inb, ok := h.Inb[u.CallbackQuery.Data]
	if ok {
		return inb.Callback(u)
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok {
		d := st.Data.(cache.AdminData)
		num := u.CallbackQuery.Data[strings.Index(u.CallbackQuery.Data, DEL)+1:]
		if isRoom(num, d.Rooms) {
			prefix := u.CallbackQuery.Data[:strings.Index(u.CallbackQuery.Data, DEL)]

			inb, ok := h.Inb[prefix]
			if ok {
				switch prefix {
				case tgBot.Text.Admin.Room2:
					d.Number = num
				case tgBot.Text.Admin.Settings.Edit.Removing4:
					d.NumberDel = num
				}
				tgBot.State.Put(cache.KeyT(u.FromChat().ID), cache.State{Is: st.Is, Data: d})
				return inb.Callback(u)
			}
		}
	}

	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *AdminHandler) Command(u *tg.Update) error {
	cmd, ok := h.Cmd[u.Message.Command()]
	if ok {
		return cmd.Action(u)
	}
	return h.Cmd[tgBot.Text.CommonCommand.Unknown].Action(u)
}

func (h *AdminHandler) Photo(u *tg.Update) error {
	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok && st.Is == tgBot.Text.Buttons.Tenant.Receipt1 {
		return h.Inp[tgBot.Text.Tenant.Receipt.Receipt2].HandleInput(u)
	} else if ok {
		return tgBot.API.SendText(u, "Сейчас мне не нужно фото.")
	} else {
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
	}
}

func (h *AdminHandler) Message(u *tg.Update) error {
	mes, ok := h.Mes[u.Message.Text]
	if ok {
		return mes.Action(u)
	}

	but, ok := h.But[u.Message.Text]
	if ok {
		return but.Action(u)
	}

	st, ok := tgBot.State.Get(cache.KeyT(u.FromChat().ID))
	if ok && st.Is != "" {
		inp, ok := h.Inp[st.Is]
		if ok {
			return inp.HandleInput(u)
		}
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

func formatNumbers(numbers []string, prefix string) (fNum [][]string, fData [][]string) {
	if (len(numbers)-1)/4 > 0 {
		fNum = make([][]string, (len(numbers)-1)/4)
		for i := range fNum {
			fNum[i] = make([]string, 4)
		}
		fNum = append(fNum, make([]string, len(numbers)%4))
	} else {
		fNum = [][]string{numbers}
	}

	for i, num := range numbers {
		//fmt.Printf("%d | %d  %s\n", i/4, i%4, num)
		fNum[i/4][i%4] = num
	}

	if (len(numbers)-1)/4 > 0 {
		fData = make([][]string, (len(numbers)-1)/4)
		for i := range fData {
			fData[i] = make([]string, 4)
		}
		fData = append(fData, make([]string, len(numbers)%4))
	} else {
		fData = append(fData, make([]string, len(numbers)))
	}
	fmt.Println(fNum, fData)
	for i := 0; i < len(fData); i++ {
		for j := 0; j < len(fData[i]); j++ {
			fData[i][j] = prefix + DEL + fNum[i][j]
		}
	}
	fmt.Println(fNum, fData)

	return fNum, fData
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
