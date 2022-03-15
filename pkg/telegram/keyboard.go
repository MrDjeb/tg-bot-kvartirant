package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard tg.ReplyKeyboardMarkup
type InKeyboard tg.InlineKeyboardMarkup

type Buttons struct {
	Tenant Tenant
	Admin  Admin
}

type Tenant struct {
	Keyboard Keyboard
	Water    InKeyboard
	Receipt  []tg.InlineKeyboardButton
}

type Admin struct {
	Keyboard Keyboard
	Rooms    InKeyboard
	Settings InKeyboard
}

func NewButtons() Buttons {
	return Buttons{
		Tenant: Tenant{
			Keyboard: Keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(tgBot.Text.Tenant.Water1), tg.NewKeyboardButton(tgBot.Text.Tenant.Receipt1)),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(tgBot.Text.Tenant.Report1)))),

			Water: InKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData((tgBot.Text.Tenant.Water.Hot_w2), (tgBot.Text.Tenant.Water.Hot_w2)),
					tg.NewInlineKeyboardButtonData(tgBot.Text.Tenant.Water.Cold_w2, tgBot.Text.Tenant.Water.Cold_w2)))),

			Receipt: tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(tgBot.Text.Tenant.Receipt.Month2, tgBot.Text.Tenant.Receipt.Month2),
				tg.NewInlineKeyboardButtonData(tgBot.Text.Tenant.Receipt.Amount2, tgBot.Text.Tenant.Receipt.Amount2),
				tg.NewInlineKeyboardButtonData(tgBot.Text.Tenant.Receipt.Receipt2, tgBot.Text.Tenant.Receipt.Receipt2),
			),
		},
		Admin: Admin{
			Keyboard: Keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(tgBot.Text.Admin.Rooms1), tg.NewKeyboardButton(tgBot.Text.Admin.Settings1)))),

			Rooms: InKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(tgBot.Text.Admin.Rooms.R2, tgBot.Text.Admin.Rooms.R2),
					tg.NewInlineKeyboardButtonData("366", "366")))),

			Settings: InKeyboard(tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(tgBot.Text.Admin.Settings.Edit2, tgBot.Text.Admin.Settings.Edit2)),
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(tgBot.Text.Admin.Settings.Contacts2, tgBot.Text.Admin.Settings.Contacts2)),
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonData(tgBot.Text.Admin.Settings.Reminder2, tgBot.Text.Admin.Settings.Reminder2)))),
		},
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
