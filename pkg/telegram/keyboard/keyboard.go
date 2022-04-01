package keyboard

import (
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

func FormatNumbers(numbers []string, prefix string) (fNum [][]string, fData [][]string) {
	if (len(numbers)-1)/4 > 0 {
		fNum = make([][]string, (len(numbers)-1)/4)
		for i := range fNum {
			fNum[i] = make([]string, 4)
		}
		fNum = append(fNum, make([]string, len(numbers)%4))
	} else {
		fNum = [][]string{numbers}
	}

	for i, num := range numbers {
		//fmt.Printf("%d | %d  %s\n", i/4, i%4, num)
		fNum[i/4][i%4] = num
	}

	if (len(numbers)-1)/4 > 0 {
		fData = make([][]string, (len(numbers)-1)/4)
		for i := range fData {
			fData[i] = make([]string, 4)
		}
		fData = append(fData, make([]string, len(numbers)%4))
	} else {
		fData = append(fData, make([]string, len(numbers)))
	}
	for i := 0; i < len(fData); i++ {
		for j := 0; j < len(fData[i]); j++ {
			fData[i][j] = prefix + DEL + fNum[i][j]
		}
	}

	return fNum, fData
}
