package telegram

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineInputer interface {
	HandleInput(u *tg.Update) error
	Callback(u *tg.Update) error
}

type TenantInlineInput struct {
	Cold_w2  InlineInputer
	Hot_w2   InlineInputer
	Month2   InlineInputer
	Amount2  InlineInputer
	Receipt2 InlineInputer
}

func (i *TenantInlineInput) New() {
	i.Cold_w2 = &Cold_w2{}
	i.Hot_w2 = &Hot_w2{}
	i.Month2 = &Month2{}
	i.Amount2 = &Amount2{}
	i.Receipt2 = &Receipt2{}
}

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

type Cold_w2 struct {
	tg.InlineKeyboardButton
	InlineInputer
}

type Hot_w2 struct {
	tg.InlineKeyboardButton
	InlineInputer
}

type Month2 struct {
	tg.InlineKeyboardButton
	InlineInputer
}
type Amount2 struct {
	tg.InlineKeyboardButton
	InlineInputer
}
type Receipt2 struct {
	tg.InlineKeyboardButton
	InlineInputer
}

func (r *Cold_w2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите показания с счётчика холодной воды. К примеру: 34,56")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	if err := tgBot.API.SendText(u, "Введите показания с счётчика холодной воды. К примеру: 34,56"); err != nil {
		return err
	}
	tgBot.State.Erase()
	tgBot.State.TenantCold_w2 = true
	return nil
}

func (r *Hot_w2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите показания с счётчика горячей воды. К примеру: 34,56")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(u.CallbackQuery.Message.Chat.ID, "Введите показания с счётчика горячей воды. К примеру: 34,56")
	if _, err := tgBot.API.Send(msg); err != nil {
		return err
	}
	tgBot.State.Erase()
	tgBot.State.TenantHot_w2 = true
	return nil
}

func (r *Month2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введите номер месяца -- число от 1 до 12.")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(u.CallbackQuery.Message.Chat.ID, "Введите номер месяца -- число от 1 до 12.")
	if _, err := tgBot.API.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[0] = 1
	case tgBot.State.TenantPayment[0] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[0] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[0] = 1
	}

	return nil
}

func (r *Amount2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(u.CallbackQuery.Message.Chat.ID, "Введиту сумму в рублях, которую вы оплатили. К примеру, 4500")
	if _, err := tgBot.API.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[1] = 1
	case tgBot.State.TenantPayment[1] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[1] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[1] = 1
	}
	return nil
}

func (r *Receipt2) Callback(u *tg.Update) error {
	callback := tg.NewCallback(u.CallbackQuery.ID, "Пришлите скрин квитанции.")
	if _, err := tgBot.API.Request(callback); err != nil {
		return err
	}
	msg := tg.NewMessage(u.CallbackQuery.Message.Chat.ID, "Пришлите скрин квитанции.")
	if _, err := tgBot.API.Send(msg); err != nil {
		return err
	}

	switch {
	case !tbool(tgBot.State.TenantPayment[0]) && !tbool(tgBot.State.TenantPayment[1]) && !tbool(tgBot.State.TenantPayment[2]):
		tgBot.State.Erase()
		tgBot.State.TenantPayment[2] = 1
	case tgBot.State.TenantPayment[2] == 2:
		tgBot.State.Erase()
		tgBot.State.TenantPayment[2] = 1
	default:
		tgBot.State.CleanProcess()
		tgBot.State.TenantPayment[2] = 1
	}
	return nil
}

func (b *Cold_w2) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(strings.Replace(u.Message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || score < 0 || score > 65.536 {
		if err := tgBot.API.SendText(u, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil //error broken
	}

	scoreDB, date := uint16(score*100), time.Now().Format(LAYOUT)
	isExist, err := tgBot.DB.Scorer.IsExistDay(u.Message.From.ID, date)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateCold_w(u.Message.From.ID, scoreDB, date); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Успешно добавлено к текущей дате!"); err != nil {
			return err
		}
	} else {
		if err := tgBot.DB.Scorer.Insert(database.Scorer{IdTg: u.Message.From.ID, Hot_w: 0, Cold_w: scoreDB, Date: date}); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Успешно создано в базе!"); err != nil {
			return err
		}
	}

	tgBot.State.Erase()

	return nil
}

func (r *Hot_w2) HandleInput(u *tg.Update) error {
	/*defer func() error {
		if err := recover(); err != nil {
			if err :=tgBot.API.SendText(u.Message.From.ID, "Мне нужен текст. К примеру: 34,56"); err != nil {
				return err
			}
		}
		return nil
	}()*/

	tidyStr := strings.TrimSpace(strings.Replace(u.Message.Text, ",", ".", 1))
	score, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || score < 0 || score > 65.536 {
		if err := tgBot.API.SendText(u, "Введите вещественное число. К примеру: 34,56"); err != nil {
			return err
		}
		return nil //error broken
	}

	scoreDB, date := uint16(score*100), time.Now().Format(LAYOUT)
	isExist, err := tgBot.DB.Scorer.IsExistDay(u.Message.From.ID, date)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateHot_w(u.Message.From.ID, scoreDB, date); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Успешно добавлено к текущей дате!"); err != nil {
			return err
		}
	} else {
		if err := tgBot.DB.Scorer.Insert(database.Scorer{IdTg: u.Message.From.ID, Hot_w: scoreDB, Cold_w: 0, Date: date}); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Успешно создано в базе!"); err != nil {
			return err
		}
	}

	tgBot.State.Erase()

	return nil
}

func (r *Month2) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	month, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || month < 1 || month > 12 {
		if err := tgBot.API.SendText(u, "Введите число от 1 до 12."); err != nil {
			return err
		}
		return nil //error broken
	}

	tgBot.State.TenantPaymentMonth = uint8(month)
	tgBot.State.TenantPayment[0] = 2

	if tgBot.State.TenantPayment[0] == 2 && tgBot.State.TenantPayment[1] == 2 && tgBot.State.TenantPayment[2] == 2 {

		if err := tgBot.DB.Payment.Insert(database.Payment{
			IdTg:      u.Message.From.ID,
			Amount:    tgBot.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(tgBot.State.TenantPaymentMonth),
			Photo:     tgBot.State.TenantPaymentReceipt}); err != nil {
			return err
		}

		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена!"); err != nil {
			return err
		}
		tgBot.State.Erase()
	} else {
		if err := tgBot.API.SendText(u, "Месяц успешно добавлен."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).Ms.Receipt1(u); err != nil {
			return err
		}
	}

	return nil
}

func (r *Amount2) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	amount, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || amount < 0 || amount > 4294967296 {
		if err := tgBot.API.SendText(u, "Введите сумму оплаты в виде числа."); err != nil {
			return err
		}
		return nil //error broken
	}

	tgBot.State.TenantPaymentAmount = uint(amount)
	tgBot.State.TenantPayment[1] = 2

	if tgBot.State.TenantPayment[0] == 2 && tgBot.State.TenantPayment[1] == 2 && tgBot.State.TenantPayment[2] == 2 {
		payment := database.Payment{
			IdTg:      u.Message.From.ID,
			Amount:    tgBot.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(tgBot.State.TenantPaymentMonth),
			Photo:     tgBot.State.TenantPaymentReceipt}

		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена c параметрами: "+fmt.Sprintf("%d ₽ | %s", payment.Amount, payment.Date)); err != nil {
			return err
		}
		tgBot.State.Erase()
	} else {
		if err := tgBot.API.SendText(u, "Cумма оплаты успешно добавлена."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).Ms.Receipt1(u); err != nil {
			return err
		}
	}

	return nil
}

func (r *Receipt2) HandleInput(u *tg.Update) error {

	fotos := u.Message.Photo

	if len(fotos) > 4 {
		if err := tgBot.API.SendText(u, "Пришлите одно фото."); err != nil {
			return err
		}
		return nil
	}
	fileURL, err := tgBot.API.GetFileDirectURL(fotos[2].FileID)
	if err != nil {
		return err
	}
	blob, err := downloadFile(fileURL)
	if err != nil {
		return err
	}

	tgBot.State.TenantPaymentReceipt = blob
	tgBot.State.TenantPayment[2] = 2

	if tgBot.State.TenantPayment[0] == 2 && tgBot.State.TenantPayment[1] == 2 && tgBot.State.TenantPayment[2] == 2 {

		if err := tgBot.DB.Payment.Insert(database.Payment{
			IdTg:      u.Message.From.ID,
			Amount:    tgBot.State.TenantPaymentAmount,
			PayMoment: time.Now().Format(LAYOUT),
			Date:      GetAverageDate(tgBot.State.TenantPaymentMonth),
			Photo:     tgBot.State.TenantPaymentReceipt}); err != nil {
			return err
		}

		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена!"); err != nil {
			return err
		}
		tgBot.State.Erase()
	} else {
		if err := tgBot.API.SendText(u, "Скрин квитанции успешно добавлен."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).Ms.Receipt1(u); err != nil {
			return err
		}
	}

	return nil
}
