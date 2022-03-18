package telegram

import (
	"fmt"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	Red map[string]RedirectResponser
	Inp map[string]InputResponser
}

type CommandResponser interface {
	Action(u *tg.Update) error
}

type MessageResponser interface {
	Action(u *tg.Update) error
}

type InputResponser interface {
	Callback(u *tg.Update) error
	HandleInput(u *tg.Update) error
}

type RedirectResponser interface {
	Callback(u *tg.Update) error
	Redirect(u *tg.Update) error
}

type ButtonResponser interface {
	Action(u *tg.Update) error
}

///////////////////
type TenantHandler struct{ HandlerResponse }

func (h *TenantHandler) New() {
	b := keyboard.NewButtons().Tenant
	h.Cmd = map[string]CommandResponser{
		tgBot.Text.CommonCommand.Start:   &TenantStart{But: b.Keyboard},
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
		return tgBot.API.SendText(u, tgBot.Text.Response.Unknown_ms)
	}
}

func (h *TenantHandler) Message(u *tg.Update) error {
	fmt.Println("MEEEEEEEEEESSSSSSSSEEEEEEEEEEEEEGGGG!!!!!!!!!")
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
	b := keyboard.NewButtons().Tenant
	h.Cmd = map[string]CommandResponser{
		tgBot.Text.CommonCommand.Start:   &TenantStart{But: b.Keyboard},
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

func (h *AdminHandler) Callback(u *tg.Update) error {
	inp, ok := h.Inp[u.CallbackQuery.Data]
	if ok {
		return inp.Callback(u)
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
	if ok {
		if st.Is == tgBot.Text.Tenant.Receipt.Receipt2 {
			return tgBot.API.SendText(u, "Пришлите фото.")
		}
		return h.Inp[st.Is].HandleInput(u)
	}
	return h.Mes[tgBot.Text.CommonCommand.Unknown].Action(u)
}
