package database

type DBTabler interface {
	Migrate() error
	Insert(r DBTabler) error
}

type TelegramID int64
type ScoreM3 uint16
type AmountRUB uint
type Date string
type Photo []byte

type Tables struct {
	Scorer  DBScorer
	Payment DBPayment
	Tenant  DBTenant
	Admin   DBAdmin
}

type Scorer struct {
	IdTg   TelegramID
	Hot_w  ScoreM3 // 0,00 - 65,536 m3
	Cold_w ScoreM3 // 0,00 - 65,536 m3
	Date   Date
}

type Payment struct {
	IdTg      TelegramID
	Amount    AmountRUB // 0 - 4294967296 Rub
	PayMoment Date
	Date      Date
	Photo     Photo
}

type Tenant struct {
	IdTg TelegramID
}

type Admin struct {
	IdTg TelegramID
}
