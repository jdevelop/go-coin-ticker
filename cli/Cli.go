package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jdevelop/go-coin-ticker/rest"
	"github.com/jdevelop/go-coin-ticker/storage"
)

func main() {

	host := flag.String("path", "", "REST path")
	addF := flag.Bool("a", false, "Add new record")
	delF := flag.Int("d", -1, "Remove record")

	flag.Parse()

	var db storage.RecordsDAO
	var err error

	if strings.HasPrefix(*host, "http") {
		db = rest.NewRestDAO(*host)
	} else {
		db, err = storage.MakeDB("/home/bofh/coins.db")
	}

	if err != nil {
		return
	}

	if *addF {
		err = add(db)
	} else if *delF != -1 {
		err = delete(db, *delF)
	} else {
		err = list(db)
	}

	if err != nil {
		log.Fatal(err)
	}

}

func delete(db storage.RecordsDAO, id int) (err error) {
	err = db.RemoveRecord(id)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func add(db storage.RecordsDAO) (err error) {
	fmt.Print("Debit code: ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	debitSym := s.Text()
	fmt.Print("Debit amount: ")
	s.Scan()
	debitAmt, err := strconv.ParseFloat(s.Text(), 64)
	if err != nil {
		return
	}
	fmt.Print("Credit code: ")
	s.Scan()
	creditSym := s.Text()
	fmt.Print("Credit amount: ")
	s.Scan()
	creditAmt, err := strconv.ParseFloat(s.Text(), 64)
	if err != nil {
		return
	}
	fmt.Print("Date (MM/DD/YYYY Hh:MM): ")
	s.Scan()
	t, err := time.Parse("01-02-2006 15:04", s.Text())
	if err != nil {
		log.Fatal(err)
	}
	err = db.AddRecord(&storage.Record{
		Date: storage.FormattedTime{Time: t},
		Credit: storage.Sale{
			Account: strings.ToLower(creditSym),
			Amount:  creditAmt,
		},
		Debit: storage.Sale{
			Account: strings.ToLower(debitSym),
			Amount:  debitAmt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record added")
	return
}

func list(db storage.RecordsDAO) (err error) {
	fmt.Println("Records content:")
	recs, err := db.GetRecords()
	if err != nil {
		return
	}

	for _, rec := range recs {
		fmt.Println(rec)
	}
	return
}
