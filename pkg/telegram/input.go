package telegram

import (
	"fmt"
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State struct {
	TenantHot_w2  bool
	TenantCold_w2 bool
}

func (s *State) Erase() {
	s.TenantHot_w2 = false
	s.TenantCold_w2 = false
}

func (b *Bot) TenantCold_w2Inp(message *tg.Message) error {

	tidyStr := strings.TrimSpace(strings.Replace(message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil {
		if err := b.handleSendText(message, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil
	}

	fmt.Println("save to db", uint16(score*100))
	if err := b.handleSendText(message, "Успешно сохранено!"); err != nil {
		return err
	}
	b.State.Erase()

	return nil
}

func (b *Bot) TenantHot_w2Inp(message *tg.Message) error {
	/*defer func() error {
		if err := recover(); err != nil {
			if err := b.handleSendText(message, "Мне нужен текст. К примеру: 34,56"); err != nil {
				return err
			}
		}
		return nil
	}()*/

	tidyStr := strings.TrimSpace(strings.Replace(message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil {
		if err := b.handleSendText(message, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil
	}

	fmt.Println("save to db", uint16(score*100))
	if err := b.handleSendText(message, "Успешно сохранено!"); err != nil {
		return err
	}
	b.State.Erase()

	return nil
}
