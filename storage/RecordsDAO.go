package storage

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	// SQLIte 3
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"
)

//PredefDebit the default account to use as a debit one.
const PredefDebit = "usd"

//PredefFee the default account for dropping the fees into.
const PredefFee = "fee"

//Unit type of the transaction
type Unit = string

//Sale defines the sale position movement.
type Sale struct {
	Account Unit    `json:"account,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

//FormattedTime is invented just for the sake of custom time format in JSON.
type FormattedTime struct {
	time.Time
}

//DatePattern the pattern for the date serialization/deserialization.
const DatePattern = "01-02-2006 15:04"

var nilTime = (time.Time{}).UnixNano()

//UnmarshalJSON converts passed bytes from string into time using the default format.
func (t *FormattedTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(DatePattern, s)
	return
}

//MarshalJSON returns the serialized version of the time object, formatted according to the format.
func (t *FormattedTime) MarshalJSON() ([]byte, error) {
	if t.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(DatePattern))), nil
}

//Record data
type Record struct {
	ID     int           `json:"id,omitempty"`
	Debit  Sale          `json:"debit,omitempty"`
	Credit Sale          `json:"credit,omitempty"`
	Date   FormattedTime `json:"date,omitempty"`
}

//RecordsDAO defines the methods to add or remove records from the underlying storage.
type RecordsDAO interface {
	AddRecord(r *Record) error
	GetRecords() ([]Record, error)
	RemoveRecord(ids ...int) error
	AggregateRecords() ([]Sale, error)
	Init() error
}

type localDB struct {
	db *sql.DB
}

func closeStmt(st ...io.Closer) {
	for _, stmt := range st {
		stmt.Close()
	}
}

func (db *localDB) RemoveRecord(ids ...int) (err error) {
	stmt, err := db.db.Prepare("DELETE FROM TRANSACTIONS WHERE ID = ?")
	defer closeStmt(stmt)
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
	stmt, err := db.db.Prepare("INSERT INTO TRANSACTIONS (DEBIT_SYM, DEBIT_AMT, CREDIT_SYM, CREDIT_AMT, TXNDATE) " +
		"VALUES (?,?,?,?,?)")
	defer closeStmt(stmt)
	if err != nil {
		return
	}
	_, err = stmt.Exec(r.Debit.Account, r.Debit.Amount, r.Credit.Account, r.Credit.Amount, r.Date.Time)
	if err != nil {
		return
	}
	return
}

func (db *localDB) GetRecords() (res []Record, err error) {
	stmt, err := db.db.Prepare("SELECT ID, DEBIT_SYM, DEBIT_AMT, CREDIT_SYM, CREDIT_AMT, TXNDATE " +
		"FROM TRANSACTIONS ORDER BY TXNDATE DESC LIMIT 10 ")
	defer closeStmt(stmt)
	if err != nil {
		return
	}
	rows, err := stmt.Query()
	defer closeStmt(rows)
	if err != nil {
		return
	}

	var (
		id                  int
		debitSym, creditSym Unit
		debitAmt, creditAmt float64
		txndate             time.Time
	)

	res = make([]Record, 0)

	for rows.Next() {
		err = rows.Scan(&id, &debitSym, &debitAmt, &creditSym, &creditAmt, &txndate)
		if err != nil {
			continue
		}
		res = append(res, Record{
			ID: id,
			Debit: Sale{
				Account: debitSym,
				Amount:  debitAmt,
			},
			Credit: Sale{
				Account: creditSym,
				Amount:  creditAmt,
			},
			Date: FormattedTime{txndate},
		})
	}

	return
}

func (db *localDB) AggregateRecords() (res []Sale, err error) {
	stmt, err := db.db.Prepare("SELECT DISTINCT(DEBIT_SYM) FROM TRANSACTIONS UNION SELECT DISTINCT(CREDIT_SYM) FROM TRANSACTIONS")
	if err != nil {
		return
	}
	defer closeStmt(stmt)

	rows, err := stmt.Query()
	defer closeStmt(rows)
	var (
		sym Unit
		amt float64
	)

	units := make([]Unit, 0)

	for rows.Next() {
		rows.Scan(&sym)
		units = append(units, sym)
	}

	if len(units) == 0 {
		return
	}

	stmt.Close()
	rows.Close()

	stmt, err = db.db.Prepare("SELECT coalesce((SELECT SUM(CREDIT_AMT) FROM TRANSACTIONS WHERE CREDIT_SYM=?),0) - " +
		"coalesce((SELECT SUM(DEBIT_AMT) FROM TRANSACTIONS WHERE DEBIT_SYM=?),0)")
	if err != nil {
		return
	}

	res = make([]Sale, 0)

	for _, unit := range units {
		rows, err := stmt.Query(unit, unit)
		if err != nil {
			log.Error("Can't query sale", err)
			continue
		}
		if rows.Next() {
			err = rows.Scan(&amt)
			if err != nil {
				log.Error("Can't scan sale data", unit, err)
				continue
			}
			res = append(res, Sale{
				Account: unit,
				Amount:  amt,
			})
		}
		rows.Close()
	}

	return

}

func (db *localDB) Init() (err error) {
	_, err = db.db.Exec(`
	CREATE TABLE TRANSACTIONS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		DEBIT_SYM VARCHAR(10),
		DEBIT_AMT DECIMAL(10,6),
		CREDIT_SYM VARCHAR(10),
		CREDIT_AMT DECIMAL(10,6),
		TXNDATE TIMESTAMP
	)`)
	if err != nil {
		return
	}
	return
}

//MakeDB creates the database for SQLite implementation.
func MakeDB(path string) (db RecordsDAO, err error) {
	_, err1 := os.Stat(path)
	init := os.IsNotExist(err1)
	_db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Error("Can't open SQLite3 driver", err)
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

//IsDebit returns true if the account the debit account used to buy a crypto?
func IsDebit(acct string) bool {
	return acct == PredefDebit
}

//IsFee returns this transaction just a fee paid for the regular transfer?
func IsFee(acct string) bool {
	return acct == PredefFee
}
