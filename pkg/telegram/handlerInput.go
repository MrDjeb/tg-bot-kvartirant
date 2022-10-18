package telegram

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/telegram/keyboard"
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
	year, m_now, _ := time.Now().Date()
	if m_now >= 7 && ((int(m_now)+7)%13) > int(m) {
		year++
	} else if m_now < 7 && (int(m_now)+6) <= int(m) {
		year--
	}
	return time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local).Format(LAYOUT[:7])
}

func getFormatCalendar(prefix string) (fNum [][]string, fData [][]string) {
	_, m_now, _ := time.Now().Date()
	fNum, fData = make([][]string, 3), make([][]string, 3)
	for i := range fNum {
		fNum[i], fData[i] = make([]string, 4), make([]string, 4)
	}
	k := uint8((m_now - 6) % 12)
	for i := range fNum {
		for j := range fNum[i] {
			fNum[i][j] = getAverageDate(k % 12)
			fData[i][j] = prefix + keyboard.DEL + fNum[i][j]
			k++
		}
	}
	return fNum, fData
}

func getFormatPayment(prefix string, payments *[]database.Payment) (fNum [][]string, fData [][]string) {
	m := make(map[string]int)
	for _, payment := range *payments {
		m[string(payment.Date)] = 0
	}
	for _, payment := range *payments {
		m[string(payment.Date)]++
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	fNum, fData = make([][]string, len(m)), make([][]string, len(m))
	for i, k := range keys {
		fNum[i], fData[i] = make([]string, 1), make([]string, 1)
		fNum[i][0] = fmt.Sprintf("`%s x%d`", k, m[k])
		fData[i][0] = prefix + keyboard.DEL + k
	}

	return fNum, fData
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

	return io.ReadAll(resp.Body)

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
		if err := tgBot.API.SendText(u, tgBot.Text.Response.Water2_inp); err != nil {
			return err
		}
		return nil //error broken
	}

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	scoreDB := database.ScoreM3(score * 100)
	var dateDB database.Date
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && d.ScoreDate != 0 {

		dateDB = database.Date(time.Now().AddDate(0, int(d.ScoreDate)-int(time.Now().Month()), 0).Format(LAYOUT))
	} else {
		dateDB = database.Date(time.Now().Format(LAYOUT))
	}

	isExist, err := tgBot.DB.Scorer.IsExistDay(num, dateDB)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateHot_w(num, scoreDB, dateDB); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Water2_change, num, dateDB)); err != nil {
			return err
		}

		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
		if !ok {
			return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
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
			if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Water2_saved, strconv.FormatFloat(float64(score.Hot_w)/100, 'f', -1, 64),
				strconv.FormatFloat(float64(score.Cold_w)/100, 'f', -1, 64), dateDB)); err != nil {
				return err
			}
			if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
				d.Erase()
				tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			}
		} else {
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			if err := tgBot.API.SendText(u, "Принял счёт за горячую воду!"); err != nil {
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
		if err := tgBot.API.SendText(u, tgBot.Text.Response.Water2_inp); err != nil {
			return err
		}
		return nil //error broken
	}

	num, err := tgBot.DB.Room.GetRoom(database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	scoreDB := database.ScoreM3(score * 100)
	var dateDB database.Date
	if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok && d.ScoreDate != 0 {
		dateDB = database.Date(time.Now().AddDate(0, int(d.ScoreDate)-int(time.Now().Month()), 0).Format(LAYOUT))
	} else {
		dateDB = database.Date(time.Now().Format(LAYOUT))
	}
	isExist, err := tgBot.DB.Scorer.IsExistDay(num, dateDB)
	if err != nil {
		return err
	}

	if isExist {
		if err := tgBot.DB.Scorer.UpdateCold_w(num, scoreDB, dateDB); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Water2_change, num, dateDB)); err != nil {
			return err
		}
		if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
			d.Erase()
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
		}
	} else {
		d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID)
		if !ok {
			return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
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
			if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Water2_saved, strconv.FormatFloat(float64(score.Hot_w)/100, 'f', -1, 64),
				strconv.FormatFloat(float64(score.Cold_w)/100, 'f', -1, 64), dateDB)); err != nil {
				return err
			}
			if d, ok := tgBot.Tenant.Cache.(*TenantCacher).Get(u.FromChat().ID); ok {
				d.Erase()
				tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			}
		} else {
			tgBot.Tenant.Cache.Put(u.FromChat().ID, d)
			if err := tgBot.API.SendText(u, "Принял счёт за холодную воду!"); err != nil {
				return err
			}
			if err := tgBot.Tenant.Handler.(*TenantHandler).HandlerResponse.But[tgBot.Text.Tenant.Water1].Action(u); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Amount2) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	amount, err := strconv.ParseFloat(tidyStr, 32)
	if err != nil || amount < 0 || amount > 4294967296 {
		if err := tgBot.API.SendText(u, tgBot.Text.Response.Amount2_inp); err != nil {
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
		return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
	}

	d.PaymentAmount = database.AmountRUB(amount)
	d.Payment[1] = true
	d.Is = ""

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentDate)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Receipt2_saved, payment.Amount, payment.Date)); err != nil {
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
		if err := tgBot.API.SendText(u, "Пришлите только одно фото."); err != nil {
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
		return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
	}

	d.PaymentReceipt = database.Photo(blob)
	d.Payment[2] = true
	d.Is = ""

	if d.Payment[0] && d.Payment[1] && d.Payment[2] {

		payment := database.Payment{
			Number:    num,
			Amount:    database.AmountRUB(d.PaymentAmount),
			PayMoment: database.Date(time.Now().Format(LAYOUT)),
			Date:      database.Date(getAverageDate(d.PaymentDate)),
			Photo:     database.Photo(d.PaymentReceipt),
		}
		if err := tgBot.DB.Payment.Insert(payment); err != nil {
			return err
		}
		if err := tgBot.API.SendText(u, fmt.Sprintf(tgBot.Text.Response.Receipt2_saved, payment.Amount, payment.Date)); err != nil {
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

	//IDs, err := tgBot.DB.Room.ReadAdmins() implement that one room can have more then one admin

	tgBot.DB.Admin.Update(database.Admin{IdTgAdmin: database.TelegramID(u.FromChat().ID), Repairer: tidyStr})

	return tgBot.API.SendText(u, fmt.Sprintf("Добавлен ремонтник %s", tidyStr))
}

func (r *AddRoom3) HandleInput(u *tg.Update) error {
	number := strings.TrimSpace(u.Message.Text)
	if len(number) > 32 {
		return tgBot.API.SendText(u, "Длина номера должна быть меньше 32 символов.")
	}

	IsExistRoom, err := tgBot.DB.Room.IsExistRoom(database.Number(number), database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	if IsExistRoom {
		return tgBot.API.SendText(u, "Данный номер комнаты занят. Пожалуйста придумайте другой:")
	}

	IsExist, err := tgBot.DB.Room.IsExist(database.Number(number), database.TelegramID(u.FromChat().ID))
	if err != nil {
		return err
	}
	if IsExist {
		if err := tgBot.API.SendText(u, "Новый квартирант будет привязан к указанной существуеющей комнате"); err != nil {
			return err
		}
	} else {
		if err := tgBot.API.SendText(u, "Создана комната с указанным номером."); err != nil {
			return err
		}
	}

	dataToken := TokenFromNum(number, u.FromChat().ID)

	d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
	}

	if d.AddingRooms == nil {
		d.AddingRooms = make(map[string]string)
	}
	d.AddingRooms[dataToken] = number
	d.Is = ""
	tgBot.Admin.Cache.Put(u.FromChat().ID, d)

	link := fmt.Sprintf(DEEPLINK, tgBot.API.Self.UserName, base64.StdEncoding.EncodeToString([]byte(dataToken)))
	if err := tgBot.API.SendText(u, fmt.Sprintf("Добавлено ожидание пользователя для комнаты под номером: %s\n%s", number, link)); err != nil {
		return err
	}
	return nil
}

func (r *Reminder2) HandleInput(u *tg.Update) error {
	msg := tg.NewMessage(u.FromChat().ID, fmt.Sprintf("Данное сообщение будет разослано: {->\n\n%s\n\n<-}", u.Message.Text))
	msg.ReplyMarkup = r.But
	tgBot.API.Send(msg)

	d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID)
	if !ok {
		return tgBot.API.SendText(u, tgBot.Text.Response.Cache_ttl)
	}
	d.Is = ""
	d.RemindText = u.Message.Text
	tgBot.Admin.Cache.Put(u.FromChat().ID, d)

	return nil
}

func (b *RemoveTenants4) HandleInput(u *tg.Update) error {

	tidyStr := strings.TrimSpace(u.Message.Text)
	IdTg, err := strconv.Atoi(tidyStr)
	if len(tidyStr) > 11 || err != nil {
		if err := tgBot.API.SendText(u, "Данный Telegram ID пользователя некорректен.\nСкопируйте ID и введите ещё раз:"); err != nil {
			return err
		}
		return nil //error broken
	}

	if _, err := tgBot.API.GetChat(tg.ChatInfoConfig{ChatConfig: tg.ChatConfig{ChatID: int64(IdTg)}}); err != nil {
		if err := tgBot.API.SendText(u, "Данный Telegram ID некорректен или не привязан к боту.\nСкопируйте ID и введите ещё раз:"); err != nil {
			return err
		}
		return nil //error broken
	}

	if d, ok := tgBot.Admin.Cache.(*AdminCacher).Get(u.FromChat().ID); ok {
		d.Is = ""
		tgBot.Admin.Cache.Put(u.FromChat().ID, d)
	}
	if err := DeleteTenant(database.TelegramID(IdTg)); err != nil {
		return err
	}

	return tgBot.API.SendText(u, fmt.Sprintf("Удалён квартирант с ID: %d", IdTg))
}
