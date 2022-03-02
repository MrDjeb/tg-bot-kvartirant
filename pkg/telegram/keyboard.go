package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type keyboard tg.ReplyKeyboardMarkup

type Buttons struct {
	Tenant
	Admin
}

type Tenant struct {
	keyboard keyboard
	Water
	Receipt
}

type Admin struct {
	keyboard keyboard
}

type Water struct {
	keyboard keyboard
}

type Receipt struct {
	add_month2   keyboard
	add_amount2  keyboard
	add_receipt2 keyboard
}

func (b *Bot) butInit() Buttons {
	return Buttons{
		Tenant: Tenant{
			keyboard: keyboard(tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Water1)), tg.NewKeyboardButton(string(b.text.Receipt1))),
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Report1))))),

			Water: Water{
				keyboard(tg.NewReplyKeyboard(
					tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Water.Hot_w2)), tg.NewKeyboardButton(string(b.text.Water.Cold_w2))))),
			},

			Receipt: Receipt{
				add_month2:   keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Receipt.Add_month2))))),
				add_amount2:  keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Receipt.Add_amount2))))),
				add_receipt2: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Receipt.Add_receipt2))))),
			},
		},
		Admin: Admin{
			keyboard: keyboard(tg.NewReplyKeyboard(tg.NewKeyboardButtonRow(tg.NewKeyboardButton(string(b.text.Admin.Rooms1))))),
		},
	}

}
