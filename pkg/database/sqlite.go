package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	ErrCountDate    = errors.New("count of date in table scorer is bad")
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
        date TEXT
    );`
	log.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBScorer) Insert(sc Scorer) error {
	log.Println("INSERT INTO scorer(idTenant, hot_w, cold_w, date) values(?,?,?,?)",
		sc.IdTg, sc.Hot_w, sc.Cold_w, sc.Date)

	_, err := r.DB.Exec("INSERT INTO scorer(idTenant, hot_w, cold_w, date) values(?,?,?,?)",
		sc.IdTg, sc.Hot_w, sc.Cold_w, sc.Date)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (s DBScorer) IsExistDay(date string) (bool, error) {
	log.Println("SELECT date FROM scorer WHERE date = ?;", date)

	rows, err := s.DB.Query("SELECT date FROM scorer WHERE date = ?;", date)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (s DBScorer) UpdateCold_w(score uint16, date string) error {
	log.Println("UPDATE scorer SET cold_w = ? WHERE date = ?", score, date)

	result, err := s.DB.Exec("UPDATE scorer SET cold_w = ? WHERE date = ?", score, date)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrCountDate
	}
	return nil
}

func (s DBScorer) UpdateHot_w(score uint16, date string) error {
	log.Println("UPDATE scorer SET hot_w = ? WHERE date = ?", score, date)

	result, err := s.DB.Exec("UPDATE scorer SET hot_w = ? WHERE date = ?", score, date)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrCountDate
	}
	return nil
}

type DBPayment struct {
	DB *sql.DB
}

func (r *DBPayment) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS payment(
        idTenant INTEGER,
        amount INTEGER,
        payMoment TEXT,
        date TEXT,
        photo BLOB
    );`
	log.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBPayment) Insert(pa Payment) error {
	log.Println("INSERT INTO payment(idTenant, amount, payMoment, date, photo) values(?,?,?,?,?)",
		pa.IdTg, pa.Amount, pa.PayMoment, pa.Date, len(pa.Photo))

	_, err := r.DB.Exec("INSERT INTO payment(idTenant, amount, payMoment, date, photo) values(?,?,?,?,?)",
		pa.IdTg, pa.Amount, pa.PayMoment, pa.Date, pa.Photo)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}

type DBTenant struct {
	DB *sql.DB
}

func (r *DBTenant) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS tenant(
        idTenant INTEGER
    );`
	log.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBTenant) IsExist(tgid int64) (bool, error) {
	log.Println("SELECT * FROM tenant WHERE idTenant = ?;", tgid)

	rows, err := r.DB.Query("SELECT * FROM tenant WHERE idTenant = ?;", tgid)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (r *DBTenant) Insert(t Tenant) error {
	log.Println("INSERT INTO tenant(idTenant) values(?)", t.IdTg)

	_, err := r.DB.Exec("INSERT INTO tenant(idTenant) values(?)", t.IdTg)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}
	return nil
}

type DBAdmin struct {
	DB *sql.DB
}

func (r *DBAdmin) Migrate() error {

	query := `
    CREATE TABLE IF NOT EXISTS admin(
        idAdmin INTEGER
    );`
	log.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBAdmin) IsExist(tgid int64) (bool, error) {
	log.Println("SELECT * FROM admin WHERE idAdmin = ?;", tgid)

	rows, err := r.DB.Query("SELECT * FROM admin WHERE idAdmin = ?;", tgid)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func Init() (Tables, error) {
	log.SetPrefix("database ")

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		return Tables{}, err
	}
	tables := Tables{DBScorer{DB: db}, DBPayment{DB: db}, DBTenant{DB: db}, DBAdmin{DB: db}}
	tables.Scorer.Migrate()
	tables.Payment.Migrate()
	tables.Tenant.Migrate()
	tables.Admin.Migrate()
	return tables, err
}
