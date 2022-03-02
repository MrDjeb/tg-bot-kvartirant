package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdStart = "start"
	msHi     = "Hi"
)

func (b *Bot) handleCommand(message *tg.Message) error {
	switch message.Command() {
	case cmdStart:
		return b.handleStartCmd(message)
	default:
		return b.handleUnknownCmd(message)
	}
}

// +++Сommands+++
func (b *Bot) handleStartCmd(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.text.Response.Start)
	msg.ReplyMarkup = b.buttons.Tenant.keyboard
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCmd(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.text.Response.Unknown_cmd)
	_, err := b.bot.Send(msg)
	return err
} // +++Сommands+++

func (b *Bot) handleMessage(message *tg.Message) error {
	switch message.Text {
	case msHi:
		return b.handleHiMs(message)
	default:
		return b.handleUnknownMs(message)
	}

}

// +++Messege+++
func (b *Bot) handleHiMs(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, fmt.Sprintf("Hello, %s!", message.From.FirstName))
	_, err := b.bot.Send(msg)
	return err
}
func (b *Bot) handleUnknownMs(message *tg.Message) error {
	msg := tg.NewMessage(message.Chat.ID, b.text.Response.Unknown_ms)
	_, err := b.bot.Send(msg)
	//fmt.Println(b.text)
	return err
} // +++Messege+++

/*

2022/03/02 19:57:39 Endpoint: getUpdates, response: {"ok":true,"result":[{"update_id":486042501,
"message":{"message_id":77,"from":{"id":657322168,"is_bot":false,"first_name":"MrDjeb","username":"MrDjeb","language_code":"ru"},
"chat":{"id":657322168,"first_name":"MrDjeb","username":"MrDjeb","type":"private"},"date":1646240259,"text":"5"}}]}

2022/03/02 19:57:42 Endpoint: getUpdates, response: {"ok":true,"result":[{"update_id":486042502,
"message":{"message_id":79,"from":{"id":657322168,"is_bot":false,"first_name":"MrDjeb","username":"MrDjeb","language_code":"ru"},
"chat":{"id":657322168,"first_name":"MrDjeb","username":"MrDjeb","type":"private"},"date":1646240262,"text":"kkk"}}]}
*/
