package api

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"os"
	"github.com/mgutz/logxi/v1"
)

type Record struct {
	Id     int       `json:"id"`
	Symbol string    `json:"symbol,omitempty"`
	Amount float64   `json:"amount,omitempty"`
	Price  float64   `json:"price,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}

type RecordsDAO interface {
	AddRecord(r *Record) error
	GetRecords() ([]Record, error)
	RemoveRecord(ids ...int) (error)
	AggregateRecords() ([]Record, error)
	Init() error
}

type localDB struct {
	db *sql.DB
}

func (db *localDB) RemoveRecord(ids ...int) (err error) {
	stmt, err := db.db.Prepare("DELETE FROM TRANSACTIONS WHERE ID = ?")
	if err != nil {
		return
	}
	for _, id := range ids {
		_, err := stmt.Exec(id)
		if err != nil {
			log.Error("Failed to remove ", id, err)
		}
	}
	return
}

func (db *localDB) AddRecord(r *Record) (err error) {
	stmt, err := db.db.Prepare("INSERT INTO TRANSACTIONS (SYMBOL, AMOUNT, PRICE, TXNDATE) VALUES (?,?,?,?)")
	if err != nil {
		return
	}
	res, err := stmt.Exec(r.Symbol, r.Amount, r.Price, r.Date)
	if err != nil {
		return
	}
	res.RowsAffected()
	return
}

func (db *localDB) GetRecords() (res []Record, err error) {
	stmt, err := db.db.Prepare("SELECT ID, SYMBOL,AMOUNT,PRICE,TXNDATE FROM TRANSACTIONS ORDER BY TXNDATE DESC LIMIT 10 ")
	if err != nil {
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var (
		id            int
		sym           string
		amount, price float64
		txndate       time.Time
	)

	res = make([]Record, 0)

	for rows.Next() {
		err = rows.Scan(&id, &sym, &amount, &price, &txndate)
		if err != nil {
			continue
		}
		res = append(res, Record{
			Id:     id,
			Symbol: sym,
			Amount: amount,
			Price:  price,
			Date:   txndate,
		})
	}

	return
}

func (db *localDB) AggregateRecords() (res []Record, err error) {
	stmt, err := db.db.Prepare("SELECT SYMBOL,SUM(AMOUNT),SUM(PRICE) FROM TRANSACTIONS GROUP BY SYMBOL")
	if err != nil {
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var (
		sym           string
		amount, price float64
	)

	res = make([]Record, 0)

	for rows.Next() {
		err = rows.Scan(&sym, &amount, &price)
		if err != nil {
			continue
		}
		res = append(res, Record{
			Symbol: sym,
			Amount: amount,
			Price:  price,
		})
	}

	return

}

func (db *localDB) Init() (err error) {
	_, err = db.db.Exec(`
	CREATE TABLE TRANSACTIONS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		SYMBOL VARCHAR(10),
		AMOUNT DECIMAL(10,6),
		PRICE DECIMAL(10,6),
		TXNDATE TIMESTAMP
	)`)
	if err != nil {
		return
	}
	return
}

func MakeDB(path string) (db RecordsDAO, err error) {
	_, err1 := os.Stat(path)
	init := os.IsNotExist(err1)
	_db, err := sql.Open("sqlite3", path)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	db = &localDB{
		db: _db,
	}

	if init {
		err = db.Init()
	}

	return
}
