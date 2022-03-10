package database

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
	IdTg   int64
	Hot_w  uint16 // 0,00 - 65,536 m3
	Cold_w uint16 // 0,00 - 65,536 m3
	Date   string
}

type Payment struct {
	IdTg      int64
	Amount    uint // 0 - 4294967296 Rub
	PayMoment string
	Date      string
	Photo     []byte
}

type Tenant struct {
	IdTg int64
}

type Admin struct {
	IdTg int64
}
