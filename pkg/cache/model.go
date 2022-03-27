package cache

import "github.com/MrDjeb/tg-bot-kvartirant/pkg/database"

type State struct {
	Is   string
	Data Data //map[string]interface[]
}

type Data interface{}

type TenantData struct {
	Payment        [3]bool
	Score          [2]bool
	PaymentMonth   uint8
	PaymentAmount  database.AmountRUB
	PaymentReceipt database.Photo
	ScoreHot_w     database.ScoreM3
	ScoreCold_w    database.ScoreM3
}

func (s *TenantData) Erase() {
	s.Score = [2]bool{false, false}
	s.Payment = [3]bool{false, false, false}
}

type AdminData struct {
	AddingRooms map[string]string
	Rooms       []string
	RoomsDel    []string
	Number      string
	NumberDel   string
}
