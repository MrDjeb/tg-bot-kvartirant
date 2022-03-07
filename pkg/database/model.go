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

/*type TelegramID int64
type ScoreM3 uint16
type AmountRUB uint*/

type Tables struct {
	Scorer  DBScorer
	Payment DBPayment
	Tenant  DBTenant
	Admin   DBAdmin
}

type Scorer struct {
	IdTenant int64
	Hot_w    uint16 // 0,00 - 65,536 m3
	Cold_w   uint16 // 0,00 - 65,536 m3
	Date     time.Time
}

type Payment struct {
	IdTenant  int64
	Amount    uint // 0 - 4294967296 Rub
	PayMoment time.Time
	Date      time.Time
	Photo     image.Image
}

type Tenant struct {
	Id int64
}

type Admin struct {
	Id int64
}
