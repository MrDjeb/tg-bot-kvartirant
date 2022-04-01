package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	ErrCountDate    = errors.New("count of date in table scorer is bad")
)

type DBScorer struct{ DB *sql.DB }

func (r *DBScorer) Migrate() error {

	query := `
    CREATE TABLE IF NOT EXISTS scorer(
        number TEXT,
        hot_w INTEGER,
        cold_w INTEGER,
        date TEXT
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBScorer) Insert(sc Scorer) error {
	logDB.Println("INSERT INTO scorer(number, hot_w, cold_w, date) values(?,?,?,?)",
		sc.Number, sc.Hot_w, sc.Cold_w, sc.Date)

	_, err := r.DB.Exec("INSERT INTO scorer(number, hot_w, cold_w, date) values(?,?,?,?)",
		sc.Number, sc.Hot_w, sc.Cold_w, sc.Date)
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

func (r *DBScorer) Delete(num Number) error {
	logDB.Println("DELETE FROM scorer WHERE number = ?", num)

	_, err := r.DB.Exec("DELETE FROM scorer WHERE number = ?", num)
	if err != nil {
		return err
	}

	return nil
}

func (s DBScorer) IsExistDay(num Number, date Date) (bool, error) {
	logDB.Println("SELECT date FROM scorer WHERE number = ? AND date = ?;", num, date)

	rows, err := s.DB.Query("SELECT date FROM scorer WHERE number = ? AND date = ?;", num, date)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (s DBScorer) UpdateCold_w(num Number, score ScoreM3, date Date) error {
	logDB.Println("UPDATE scorer SET cold_w = ? WHERE number = ? AND date = ?", score, num, date)

	result, err := s.DB.Exec("UPDATE scorer SET cold_w = ? WHERE number = ? AND date = ?", score, num, date)
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

func (s DBScorer) UpdateHot_w(num Number, score ScoreM3, date Date) error {
	logDB.Println("UPDATE scorer SET hot_w = ? WHERE number = ? AND date = ?", score, num, date)

	result, err := s.DB.Exec("UPDATE scorer SET hot_w = ? WHERE number = ? AND date = ?", score, num, date)
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

func (r *DBScorer) Read(num Number) ([]Scorer, error) {
	logDB.Println("SELECT * FROM scorer WHERE number = ?;", num)

	rows, err := r.DB.Query("SELECT * FROM scorer WHERE number = ?;", num)
	if err != nil {
		return []Scorer{}, err
	}
	defer rows.Close()
	var scorers []Scorer
	for rows.Next() {
		var p Scorer
		if err := rows.Scan(&p.Number, &p.Hot_w, &p.Cold_w, &p.Date); err != nil {
			return []Scorer{}, err
		}
		scorers = append(scorers, p)
	}
	return scorers, nil
}

type DBPayment struct{ DB *sql.DB }

func (r *DBPayment) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS payment(
        number TEXT,
        amount INTEGER,
        payMoment TEXT,
        date TEXT,
        photo BLOB
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBPayment) Insert(pa Payment) error {
	logDB.Println("INSERT INTO payment(number, amount, payMoment, date, photo) values(?,?,?,?,?)",
		pa.Number, pa.Amount, pa.PayMoment, pa.Date, len(pa.Photo))

	_, err := r.DB.Exec("INSERT INTO payment(number, amount, payMoment, date, photo) values(?,?,?,?,?)",
		pa.Number, pa.Amount, pa.PayMoment, pa.Date, pa.Photo)
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

func (r *DBPayment) Delete(num Number) error {
	logDB.Println("DELETE FROM payment WHERE number = ?", num)

	_, err := r.DB.Exec("DELETE FROM payment WHERE number = ?", num)
	if err != nil {
		return err
	}

	return nil
}

func (r *DBPayment) Read(num Number) ([]Payment, error) {
	logDB.Println("SELECT * FROM payment WHERE number = ?;", num)

	rows, err := r.DB.Query("SELECT * FROM payment WHERE number = ?;", num)
	if err != nil {
		return []Payment{}, err
	}
	defer rows.Close()
	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.Number, &p.Amount, &p.PayMoment, &p.Date, &p.Photo); err != nil {
			return []Payment{}, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

type DBTenant struct{ DB *sql.DB }

func (r *DBTenant) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS tenant(
        idTenant INTEGER
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBTenant) IsExist(tgid TelegramID) (bool, error) {
	logDB.Println("SELECT * FROM tenant WHERE idTenant = ?;", tgid)

	rows, err := r.DB.Query("SELECT * FROM tenant WHERE idTenant = ?;", tgid)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (r *DBTenant) Insert(t Tenant) error {
	logDB.Println("INSERT INTO tenant(idTenant) values(?)", t.IdTg)

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

type DBAdmin struct{ DB *sql.DB }

func (r *DBAdmin) Migrate() error {

	query := `
    CREATE TABLE IF NOT EXISTS admin(
        idAdmin INTEGER,
		repairer TEXT
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBAdmin) Insert(a Admin) error {
	logDB.Println("INSERT INTO admin(idAdmin, repairer) values(?,?)", a.IdTgAdmin, a.Repairer)

	_, err := r.DB.Exec("INSERT INTO admin(idAdmin, repairer) values(?,?)", a.IdTgAdmin, a.Repairer)
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

func (r *DBAdmin) GetRepairer(tgid TelegramID) (string, error) {
	logDB.Println("SELECT repairer FROM admin WHERE idAdmin = ?;", tgid)

	row := r.DB.QueryRow("SELECT repairer FROM admin WHERE idAdmin = ?;", tgid)
	var username string
	return username, row.Scan(&username)
}

func (r *DBAdmin) IsExist(tgid TelegramID) (bool, error) {
	logDB.Println("SELECT * FROM admin WHERE idAdmin = ?;", tgid)

	rows, err := r.DB.Query("SELECT * FROM admin WHERE idAdmin = ?;", tgid)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (s DBAdmin) Update(a Admin) error {
	logDB.Println("UPDATE admin SET repairer = ? WHERE idAdmin = ?", a.Repairer, a.IdTgAdmin)

	result, err := s.DB.Exec("UPDATE admin SET repairer = ? WHERE idAdmin = ?", a.Repairer, a.IdTgAdmin)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrUpdateFailed
	}
	return nil
}

type DBRoom struct{ DB *sql.DB }

func (r *DBRoom) Migrate() error {

	query := `
    CREATE TABLE IF NOT EXISTS room(
		idAdmin INTEGER,
        idTenant INTEGER,
        number TEXT
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBRoom) Insert(o Room) error {
	logDB.Println("INSERT INTO room(idAdmin, idTenant, number) values(?,?,?)", o.IdTgAdmin, o.IdTgTenant, o.Number)

	_, err := r.DB.Exec("INSERT INTO room(idAdmin, idTenant, number) values(?,?,?)", o.IdTgAdmin, o.IdTgTenant, o.Number)
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

func (r *DBRoom) Delete(num Number) error {
	logDB.Println("DELETE FROM room WHERE number = ?", num)

	res, err := r.DB.Exec("DELETE FROM room WHERE number = ?", num)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

func (r *DBRoom) Read(tgid TelegramID) ([]Room, error) {
	logDB.Println("SELECT * FROM room WHERE idAdmin = ?;", tgid)

	rows, err := r.DB.Query("SELECT * FROM room WHERE idAdmin = ?;", tgid)
	if err != nil {
		return []Room{}, err
	}
	defer rows.Close()
	var rooms []Room
	for rows.Next() {
		var p Room
		if err := rows.Scan(&p.IdTgAdmin, &p.IdTgTenant, &p.Number); err != nil {
			return []Room{}, err
		}
		rooms = append(rooms, p)
	}
	return rooms, nil
}

func (r *DBRoom) ReadTenants(num Number) ([]Room, error) {
	logDB.Println("SELECT * FROM room WHERE number = ?;", num)

	rows, err := r.DB.Query("SELECT * FROM room WHERE number = ?;", num)
	if err != nil {
		return []Room{}, err
	}
	defer rows.Close()
	var rooms []Room
	for rows.Next() {
		var p Room
		if err := rows.Scan(&p.IdTgAdmin, &p.IdTgTenant, &p.Number); err != nil {
			return []Room{}, err
		}
		rooms = append(rooms, p)
	}
	return rooms, nil
}

func (r *DBRoom) GetRoom(tgid TelegramID) (Number, error) {
	logDB.Println("SELECT number FROM room WHERE idTenant = ?;", tgid)

	row := r.DB.QueryRow("SELECT number FROM room WHERE idTenant = ?;", tgid)
	var num Number
	return num, row.Scan(&num)
}

func (r *DBRoom) GetAdmin(tgid TelegramID) (TelegramID, error) {
	logDB.Println("SELECT idAdmin FROM room WHERE idTenant = ?;", tgid)

	row := r.DB.QueryRow("SELECT idAdmin FROM room WHERE idTenant = ?;", tgid)
	var id TelegramID
	return id, row.Scan(&id)
}

var logDB *log.Logger

func Init() (Tables, error) {
	logDB = log.New(os.Stderr, "[SQLITE] ", log.LstdFlags|log.Lmsgprefix)

	db, err := sql.Open("sqlite3", "stage.db")
	if err != nil {
		return Tables{}, err
	}
	tables := Tables{DBScorer{DB: db}, DBPayment{DB: db}, DBTenant{DB: db}, DBAdmin{DB: db}, DBRoom{DB: db}}
	tables.Scorer.Migrate()
	tables.Payment.Migrate()
	tables.Tenant.Migrate()
	tables.Admin.Migrate()
	tables.Room.Migrate()
	return tables, err
}
