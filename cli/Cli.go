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

	"github.com/jdevelop/go-coin-ticker/cointicker"
)

func main() {

	addF := flag.Bool("a", false, "Add new record")
	delF := flag.Int("d", -1, "Remove record")

	flag.Parse()

	db, err := cointicker.MakeDB("/home/bofh/coins.db")
	if err != nil {
		return
	}

	if *addF {
		add(db)
	} else if *delF != -1 {
		delete(db, *delF)
	} else {
		list(db)
	}

}

func delete(db cointicker.RecordsDAO, id int) {
	err := db.RemoveRecord(id)
	if err != nil {
		log.Fatal(err)
	}
}

func add(db cointicker.RecordsDAO) {
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
	fmt.Print("Date (MM-DD-YYYY Hh:MM): ")
	s.Scan()
	t, err := time.Parse("01/02/2006 15:04", s.Text())
	if err != nil {
		log.Fatal(err)
	}
	err = db.AddRecord(&cointicker.Record{
		Date: t,
		Credit: cointicker.Sale{
			Account: strings.ToLower(creditSym),
			Amount:  creditAmt,
		},
		Debit: cointicker.Sale{
			Account: strings.ToLower(debitSym),
			Amount:  debitAmt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record added")
}

func list(db cointicker.RecordsDAO) {
	fmt.Println("Database content:")
	recs, err := db.GetRecords()
	if err != nil {
		return
	}

	for _, rec := range recs {
		fmt.Println(rec)
	}

}
