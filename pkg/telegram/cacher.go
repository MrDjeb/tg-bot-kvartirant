package telegram

import (
	"log"

	"github.com/MrDjeb/tg-bot-kvartirant/pkg/cache"
	"github.com/MrDjeb/tg-bot-kvartirant/pkg/database"
)

type TenantData struct {
	Is             string
	Payment        [3]bool
	Score          [2]bool
	PaymentDate    uint8
	PaymentAmount  database.AmountRUB
	PaymentReceipt database.Photo
	ScoreHot_w     database.ScoreM3
	ScoreCold_w    database.ScoreM3
	ScoreDate      uint8
}

func (s *TenantData) Erase() {
	s.Score = [2]bool{false, false}
	s.Payment = [3]bool{false, false, false}
	s.ScoreDate = 0
	s.Is = ""
}

type AdminData struct {
	Is          string
	AddingRooms map[string]string
	Rooms       []string
	RoomsDel    []string
	Number      string
	NumberDel   string
	ShowPayment string
	RemindText  string
}

type Cacher interface {
	New()
	Put(key int64, val interface{})
}

type TenantCacher struct {
	Get func(key int64) (val TenantData, ok bool)
}

func (h *TenantCacher) New() {
	h.Get = func(key int64) (val TenantData, ok bool) {
		switch value, ok := tgBot.Cache.Get(cache.KeyT(key)); value.(type) {
		case nil:
			return TenantData{}, ok
		case TenantData:
			return value.(TenantData), ok
		default:
			log.Printf("Getting wrong value type by id: %d (need TenantData)\n", key)
			return TenantData{}, ok
		}
	}
}

func (h *TenantCacher) Put(key int64, val interface{}) { tgBot.Cache.Put(cache.KeyT(key), val) }

type AdminCacher struct {
	Get func(key int64) (val AdminData, ok bool)
}

func (h *AdminCacher) New() {
	h.Get = func(key int64) (val AdminData, ok bool) {
		switch value, ok := tgBot.Cache.Get(cache.KeyT(key)); value.(type) {
		case nil:
			return AdminData{}, ok
		case AdminData:
			return value.(AdminData), ok
		default:
			log.Printf("Getting wrong value type by id: %d (need AdminData)\n", key)
			return AdminData{}, ok
		}
	}

}

func (h *AdminCacher) Put(key int64, val interface{}) { tgBot.Cache.Put(cache.KeyT(key), val) }

type UnknownCacher struct {
	Get func(key int64) (val interface{}, ok bool)
}

func (h *UnknownCacher) New() {
	h.Get = func(key int64) (val interface{}, ok bool) {
		value, ok := tgBot.Cache.Get(cache.KeyT(key))
		return value, ok
	}
}
func (h *UnknownCacher) Put(key int64, val interface{}) { tgBot.Cache.Put(cache.KeyT(key), val) }
