package database

import (
	"image"
	"time"
)

type DBTable interface {
	Migrate() error
	Insert(r DBTable) (*DBTable, error)
	All() ([]DBTable, error)
	GetByName(name string) (*DBTable, error)
	Update(id int64, updated DBTable) (*DBTable, error)
	Delete(id int64) error
}

type TelegramID int64
type ScoreM3 uint16
type AmountRub uint

type Scorer struct {
	IdTenant TelegramID
	Hot_w    ScoreM3 // 0,00 - 65,536 m3
	Cold_w   ScoreM3 // 0,00 - 65,536 m3
	Date     time.Time
}

type Payment struct {
	IdTenant  TelegramID
	Amount    AmountRub // 0 - 4294967296 Rub
	PayMoment time.Time
	Date      time.Time
	Photo     image.Image
}
