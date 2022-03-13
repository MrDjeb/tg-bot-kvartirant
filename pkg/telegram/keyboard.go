package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type keyboard tg.ReplyKeyboardMarkup
type inKeyboard tg.InlineKeyboardMarkup

type Buttons struct {
	Tenant
	Admin keyboard
}

type Tenant struct {
	keyboard keyboard
	Water    inKeyboard
	Receipt  []tg.InlineKeyboardButton
}

func NewButtons() Buttons {
	return Buttons{
		Tenant: Tenant{
			keyboard: keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(tgBot.Text.Water1)), tg.NewKeyboardButton(string(tgBot.Text.Receipt1))),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(tgBot.Text.Report1))))),

			Water: inKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(string(tgBot.Text.Water.Hot_w2), string(tgBot.Text.Water.Hot_w2)), tg.NewInlineKeyboardButtonData(string(tgBot.Text.Water.Cold_w2), string(tgBot.Text.Water.Cold_w2))))),

			Receipt: tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(string(tgBot.Text.Receipt.Add_month2), string(tgBot.Text.Receipt.Add_month2)),
				tg.NewInlineKeyboardButtonData(string(tgBot.Text.Receipt.Add_amount2), string(tgBot.Text.Receipt.Add_amount2)),
				tg.NewInlineKeyboardButtonData(string(tgBot.Text.Receipt.Add_receipt2), string(tgBot.Text.Receipt.Add_receipt2)),
			),
		},
		Admin: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(tgBot.Text.Admin.Rooms1))))),
	}
}

type State struct {
	TenantHot_w2         bool
	TenantCold_w2        bool
	TenantPayment        [3]int
	TenantPaymentMonth   uint8
	TenantPaymentAmount  uint
	TenantPaymentReceipt []byte
}

func (s *State) Erase() {
	s.TenantHot_w2 = false
	s.TenantCold_w2 = false
	s.TenantPayment = [3]int{0, 0, 0} // 0 - isn't, 1 - processing, 2 - done
}

func (s *State) CleanProcess() {
	s.TenantHot_w2 = false
	s.TenantCold_w2 = false
	for i := range s.TenantPayment {
		if s.TenantPayment[i] == 1 {
			s.TenantPayment[i] = 0
		}
	}
}
