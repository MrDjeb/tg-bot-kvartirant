package database

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type DBScorer struct {
	DB *sql.DB
}

func (r *DBScorer) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS scorer(
        idTenant INTEGER,
        hot_w INTEGER,
        cold_w INTEGER,
        date INTEGER
    );`

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBScorer) Insert(sc Scorer) (*Scorer, error) {
	res, err := r.DB.Exec("INSERT INTO scorer(idTenant, hot_w, cold_w, date) values(?,?,?,?)",
		sc.IdTenant, sc.Hot_w, sc.Cold_w, sc.Date)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	sc.IdTenant = id

	return &sc, nil
}

type DBPayment struct {
	DB *sql.DB
}

func (r *DBPayment) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS payment(
        idTenant INTEGER,
        amount INTEGER,
        payMoment INTEGER,
        date INTEGER,
        photo BLOB
    );`

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBPayment) Insert(pa Payment) (*Payment, error) {
	res, err := r.DB.Exec("INSERT INTO payment(idTenant, amount, payMoment, date, photo) values(?,?,?,?,?)",
		pa.IdTenant, pa.Amount, pa.PayMoment, pa.Date, pa.Photo)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	pa.IdTenant = id

	return &pa, nil
}

type DBTenant struct {
	DB *sql.DB
}

func (r *DBTenant) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS tenant(
        idTenant INTEGER
    );`

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBTenant) IsExist(tgid int64) bool {

	return true
}

type DBAdmin struct {
	DB *sql.DB
}

func (r *DBAdmin) IsExist(tgid int64) bool {
	return false

}

func Init() (Tables, error) {
	db, err := sql.Open("sqlite3", "bot.db")
	if err != nil {
		return Tables{}, err
	}
	return Tables{DBScorer{DB: db}, DBPayment{DB: db}, DBTenant{DB: db}, DBAdmin{DB: db}}, err
}
