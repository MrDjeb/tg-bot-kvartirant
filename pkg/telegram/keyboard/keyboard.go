package keyboard

import (
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard tg.ReplyKeyboardMarkup
type InKeyboard tg.InlineKeyboardMarkup

const DEL = "$"

func MakeKeyboard(names ...[]string) Keyboard {
	var buttons [][]tg.KeyboardButton
	for _, row := range names {
		var butRow []tg.KeyboardButton
		for _, name := range row {
			butRow = append(butRow, tg.NewKeyboardButton(name))
		}
		buttons = append(buttons, tg.NewKeyboardButtonRow(butRow...))

	}
	return Keyboard(tg.NewReplyKeyboard(buttons...))
}

func MakeInKeyboard(names [][]string, data [][]string) InKeyboard {
	if len(names) != len(data) {
		return InKeyboard{}
	}
	for i := 0; i < len(names); i++ {
		if len(names[i]) != len(data[i]) {
			return InKeyboard{}
		}
	}

	var buttons [][]tg.InlineKeyboardButton
	for i := 0; i < len(names); i++ {
		var butRow []tg.InlineKeyboardButton
		for j := 0; j < len(names[i]); j++ {
			butRow = append(butRow, tg.NewInlineKeyboardButtonData(names[i][j], data[i][j]))
		}
		buttons = append(buttons, tg.NewInlineKeyboardRow(butRow...))

	}
	return InKeyboard(tg.NewInlineKeyboardMarkup(buttons...))
}

func FormatNumbers(mapNum map[string]int, prefix string) (fNum [][]string, fData [][]string) {
	n := len(mapNum)
	if (n-1)/4 > 0 {
		fNum = make([][]string, (n-1)/4)
		fData = make([][]string, (n-1)/4)
		for i := range fNum {
			fNum[i] = make([]string, 4)
			fData[i] = make([]string, 4)
		}
		fNum = append(fNum, make([]string, n%4))
		fData = append(fData, make([]string, n%4))
	} else {
		fNum = append(fNum, make([]string, n))
		fData = append(fData, make([]string, n))
	}

	i := 0
	for k, v := range mapNum {
		fNum[i/4][i%4] = k + "✖" + strconv.Itoa(v)
		fData[i/4][i%4] = prefix + DEL + k
		i++
	}

	return fNum, fData
}

func MakeFormatMonth(prefix string) (fNum [][]string, fData [][]string) {
	fNum, fData = make([][]string, 3), make([][]string, 3)
	for i := range fNum {
		fNum[i], fData[i] = make([]string, 4), make([]string, 4)
	}

	for i := 0; i < 9; i++ {
		fNum[i/4][i%4] = string('1' + rune(i))
		fData[i/4][i%4] = prefix + DEL + fNum[i/4][i%4]
	}
	fNum[2][1] = "10"
	fNum[2][2] = "11"
	fNum[2][3] = "12"

	for i := range fNum {
		for j := range fNum[i] {
			fData[i][j] = prefix + DEL + fNum[i][j]
		}
	}
	return fNum, fData
}
