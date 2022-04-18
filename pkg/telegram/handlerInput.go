package telegram

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	LAYOUT   = "2006-01-02"
	DEEPLINK = "https://t.me/%s?start=%s"
	//EXPIRE_TIME = 60 //minuts
)

/*
func tbool(n int) bool {
	return !(n == 0)
}*/

func tint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TokenFromNum(number string, idAdmin int64) string {
	/*Claims := jwt.MapClaims{}
	Claims["number"] = number
	Claims["idAdmin"] = idAdmin
	Claims["exp"] = time.Now().Add(time.Minute * EXPIRE_TIME).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString([]byte(tgBot.Text.JwtKey))*/
	h := md5.New()
	io.WriteString(h, number)
	return fmt.Sprintf("%x", h.Sum(nil)) + strconv.FormatInt(idAdmin, 10)
}

func getAverageDate(m uint8) string {
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

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	scoreDB := database.ScoreM3(score * 100)
	dateDB := database.Date(time.Now().Format(LAYOUT))
	isExist, err := tgBot.DB.Scorer.IsExistDay(num, dateDB)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateHot_w(num, scoreDB, dateDB); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf("Показания горячей воды в 〈%s〉 за сегодня успешно изменены.", num)); err != nil {
			return err
		}

		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
		if !ok {
			return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
		}

		d.ScoreHot_w = scoreDB
		d.Score[0] = true

		if d.Score[0] && d.Score[1] {

			score := database.Scorer{
				Number: num,
				Hot_w:  d.ScoreHot_w,
				Cold_w: d.ScoreCold_w,
				Date:   dateDB,
			}
			if err := tgBot.DB.Scorer.Insert(score); err != nil {
				return err
			}
			if err := tgBot.API.SendText(u, "Счёт за воду успешно сохранен c параметрами: "+fmt.Sprintf("%s Hot m3 | %s Cold m3", strconv.FormatFloat(float64(score.Hot_w)/100, 'f', -1, 64), strconv.FormatFloat(float64(score.Cold_w)/100, 'f', -1, 64))); err != nil {
				return err
			}
			if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
				d.Erase()
				tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			}
		} else {
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			if err := tgBot.API.SendText(u, "Счёт за горячую воду внесён."); err != nil {
				return err
			}
			if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Water1].Action(u); err != nil {
				return err
			}
		}
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

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	scoreDB := database.ScoreM3(score * 100)
	dateDB := database.Date(time.Now().Format(LAYOUT))
	isExist, err := tgBot.DB.Scorer.IsExistDay(num, dateDB)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateCold_w(num, scoreDB, dateDB); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf("Показания холодной воды в 〈%s〉 за сегодня успешно изменены.", num)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
		if !ok {
			return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
		}

		d.ScoreCold_w = scoreDB
		d.Score[1] = true

		if d.Score[0] && d.Score[1] {

			score := database.Scorer{
				Number: num,
				Hot_w:  d.ScoreHot_w,
				Cold_w: d.ScoreCold_w,
				Date:   dateDB,
			}
			if err := tgBot.DB.Scorer.Insert(score); err != nil {
				return err
			}
			if err := tgBot.API.SendText(u, "Счёт за воду успешно сохранен c параметрами: "+fmt.Sprintf("%s Hot m3 | %s Cold m3", strconv.FormatFloat(float64(score.Hot_w)/100, 'f', -1, 64), strconv.FormatFloat(float64(score.Cold_w)/100, 'f', -1, 64))); err != nil {
				return err
			}
			if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
				d.Erase()
				tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			}
		} else {
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			if err := tgBot.API.SendText(u, "Счёт за холодную воду внесён."); err != nil {
				return err
			}
			if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Water1].Action(u); err != nil {
				return err
			}
		}
	}

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

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}

	d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
	}

	d.PaymentMonth = uint8(month)
	d.Payment[0] = true

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentMonth)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена c параметрами: "+fmt.Sprintf("%d ₽ | %s", payment.Amount, payment.Date)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		if err := tgBot.API.SendText(u, "Месяц успешно добавлен."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Receipt1].Action(u); err != nil {
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

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}

	d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
	}

	d.PaymentAmount = database.AmountRUB(amount)
	d.Payment[1] = true

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentMonth)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена c параметрами: "+fmt.Sprintf("%d ₽ | %s", payment.Amount, payment.Date)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		if err := tgBot.API.SendText(u, "Сумма оплаты добавлена."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Receipt1].Action(u); err != nil {
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
	//fmt.Println(tg.FileID(fotos[2].FileID).UploadData()) Затычка в библиотеке, функция просто вызывает panic
	fileURL, err := tgBot.API.GetFileDirectURL(fotos[2].FileID)
	if err != nil {
		return err
	}
	blob, err := downloadFile(fileURL)
	if err != nil {
		return err
	}

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}

	d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
	}

	d.PaymentReceipt = database.Photo(blob)
	d.Payment[2] = true

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentMonth)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, "Квитанция успешно сохранена c параметрами: "+fmt.Sprintf("%d ₽ | %s", payment.Amount, payment.Date)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		if err := tgBot.API.SendText(u, "Скрин оплаты добавлен."); err != nil {
			return err
		}
		if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Receipt1].Action(u); err != nil {
			return err
		}
	}

	return nil
}

/////////////////////////////ADMIN/////////////////////////////

func (b *Contacts2) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	if !strings.HasPrefix(tidyStr, "@") {
		if err := tgBot.API.SendText(u, "Введите username"); err != nil {
			return err
		}
		return nil //error broken
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = ""
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	}

	tgBot.DB.Admin.Update(database.Admin{IdTgAdmin: database.TelegramID(u.FromChat().ID), Repairer: tidyStr})

	return tgBot.API.SendText(u, fmt.Sprintf("Добавлен ремонтник %s", tidyStr))
}

func (r *AddRoom3) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	/*number, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || number < 0 || number > 4294967296 {
		if err := tgBot.API.SendText(u, "Введите номер комнаты в виде числа."); err != nil {
			return err
		}
		return nil //error broken
	}*/
	dataToken := TokenFromNum(tidyStr, u.FromChat().ID)

	d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, "Пожалуйства, начните операцию ввода заного")
	}

	if d.AddingRooms == nil {
		d.AddingRooms = make(map[string]string)
	}
	d.AddingRooms[dataToken] = tidyStr
	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = ""
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	}

	link := fmt.Sprintf(DEEPLINK, tgBot.API.Self.UserName, base64.StdEncoding.EncodeToString([]byte(dataToken)))
	if err := tgBot.API.SendText(u, fmt.Sprintf("Добавлено ожидание пользователя для квартиры под номером: %s\n%s", tidyStr, link)); err != nil {
		return err
	}
	return nil
}
