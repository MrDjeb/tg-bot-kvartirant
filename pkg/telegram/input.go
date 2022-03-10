package telegram

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const LAYOUT = "2006-01-02"

func tbool(n int) bool {
	return !(n == 0)
}

func GetAverageDate(m uint8) string {
	tmp := time.Now()
	year, m_now, _ := tmp.Date()
	if m_now >= 7 && ((int(m_now)+7)/13) > int(m) {
		year++
	}
	if m_now < 7 && (int(m_now)+6) <= int(m) {
		year--
	}
	tmp = time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
	return tmp.Format(LAYOUT[:7])
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

func (b *Bot) TenantCold_w2Inp(message *tg.Message) error {

	tidyStr := strings.TrimSpace(strings.Replace(message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || score < 0 || score > 65.536 {
		if err := b.handleSendText(message.From.ID, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil //error broken
	}

	scoreDB, date := uint16(score*100), time.Now().Format(LAYOUT)
	isExist, err := b.DB.Scorer.IsExistDay(date)
	if err != nil {
		return err
	}

	if isExist {
		if err := b.DB.Scorer.UpdateCold_w(scoreDB, date); err != nil {
			return err
		}
		if err := b.handleSendText(message.From.ID, "Успешно добавлено к текущей дате!"); err != nil {
			return err
		}
	} else {
		if err := b.DB.Scorer.Insert(database.Scorer{IdTg: message.From.ID, Hot_w: 0, Cold_w: scoreDB, Date: date}); err != nil {
			return err
		}
		if err := b.handleSendText(message.From.ID, "Успешно создано в базе!"); err != nil {
			return err
		}
	}

	b.State.Erase()

	return nil
}

func (b *Bot) TenantHot_w2Inp(message *tg.Message) error {
	/*defer func() error {
		if err := recover(); err != nil {
			if err := b.handleSendText(message.From.ID, "Мне нужен текст. К примеру: 34,56"); err != nil {
				return err
			}
		}
		return nil
	}()*/

	tidyStr := strings.TrimSpace(strings.Replace(message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || score < 0 || score > 65.536 {
		if err := b.handleSendText(message.From.ID, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil //error broken
	}

	scoreDB, date := uint16(score*100), time.Now().Format(LAYOUT)
	isExist, err := b.DB.Scorer.IsExistDay(date)
	if err != nil {
		return err
	}

	if isExist {
		if err := b.DB.Scorer.UpdateHot_w(scoreDB, date); err != nil {
			return err
		}
		if err := b.handleSendText(message.From.ID, "Успешно добавлено к текущей дате!"); err != nil {
			return err
		}
	} else {
		if err := b.DB.Scorer.Insert(database.Scorer{IdTg: message.From.ID, Hot_w: scoreDB, Cold_w: 0, Date: date}); err != nil {
			return err
		}
		if err := b.handleSendText(message.From.ID, "Успешно создано в базе!"); err != nil {
			return err
		}
	}

	b.State.Erase()

	return nil
}

func (b *Bot) TenantAdd_month2Inp(message *tg.Message) error {

	tidyStr := strings.TrimSpace(message.Text)
	month, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || month < 1 || month > 12 {
		if err := b.handleSendText(message.From.ID, "Введите число от 1 до 12."); err != nil {
			return err
		}
		return nil //error broken
	}

	b.State.TenantPaymentMonth = uint8(month)
	b.State.TenantPayment[0] = 2

	if b.State.TenantPayment[0] == 2 && b.State.TenantPayment[1] == 2 && b.State.TenantPayment[2] == 2 {

		if err := b.DB.Payment.Insert(database.Payment{
			IdTg:      message.From.ID,
			Amount:    b.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(b.State.TenantPaymentMonth),
			Photo:     b.State.TenantPaymentReceipt}); err != nil {
			return err
		}

		if err := b.handleSendText(message.From.ID, "Квитанция успешно сохранена!"); err != nil {
			return err
		}
		b.State.Erase()
	} else {
		if err := b.handleSendText(message.From.ID, "Месяц успешно добавлен."); err != nil {
			return err
		}
		if err := b.TenantReceipt1Ms(message); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) TenantAdd_amount2Inp(message *tg.Message) error {

	tidyStr := strings.TrimSpace(message.Text)
	amount, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || amount < 0 || amount > 4294967296 {
		if err := b.handleSendText(message.From.ID, "Введите сумму оплаты в виде числа."); err != nil {
			return err
		}
		return nil //error broken
	}

	b.State.TenantPaymentAmount = uint(amount)
	b.State.TenantPayment[1] = 2

	if b.State.TenantPayment[0] == 2 && b.State.TenantPayment[1] == 2 && b.State.TenantPayment[2] == 2 {

		if err := b.DB.Payment.Insert(database.Payment{
			IdTg:      message.From.ID,
			Amount:    b.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(b.State.TenantPaymentMonth),
			Photo:     b.State.TenantPaymentReceipt}); err != nil {
			return err
		}

		if err := b.handleSendText(message.From.ID, "Квитанция успешно сохранена!"); err != nil {
			return err
		}
		b.State.Erase()
	} else {
		if err := b.handleSendText(message.From.ID, "Cумма оплаты успешно добавлена."); err != nil {
			return err
		}
		if err := b.TenantReceipt1Ms(message); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) TenantAdd_receipt2Inp(message *tg.Message) error {

	fotos := message.Photo

	if len(fotos) > 4 {
		if err := b.handleSendText(message.From.ID, "Пришлите одно фото."); err != nil {
			return err
		}
		return nil
	}
	fileURL, err := b.Api.GetFileDirectURL(fotos[2].FileID)
	if err != nil {
		return err
	}
	blob, err := downloadFile(fileURL)
	if err != nil {
		return err
	}

	b.State.TenantPaymentReceipt = blob
	b.State.TenantPayment[2] = 2

	if b.State.TenantPayment[0] == 2 && b.State.TenantPayment[1] == 2 && b.State.TenantPayment[2] == 2 {

		if err := b.DB.Payment.Insert(database.Payment{
			IdTg:      message.From.ID,
			Amount:    b.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(b.State.TenantPaymentMonth),
			Photo:     b.State.TenantPaymentReceipt}); err != nil {
			return err
		}

		if err := b.handleSendText(message.From.ID, "Квитанция успешно сохранена!"); err != nil {
			return err
		}
		b.State.Erase()
	} else {
		if err := b.handleSendText(message.From.ID, "Скрин квитанции успешно добавлен."); err != nil {
			return err
		}
		if err := b.TenantReceipt1Ms(message); err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	return ioutil.ReadAll(resp.Body)

}
