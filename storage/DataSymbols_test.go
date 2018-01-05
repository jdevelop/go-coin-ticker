package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TmpPath = "/tmp/coins.db.test"

func Test_DBAccess(t *testing.T) {

	os.Remove(TmpPath)

	db, err := MakeDB(TmpPath)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	recs := []Record{
		{
			Credit: Sale{
				Account: "XRP",
				Amount:  100.0,
			},
			Debit: Sale{
				Account: "USD",
				Amount:  100.0,
			},
			Date: FormattedTime{time.Now().Add(1 * time.Hour)},
		},
		{
			Credit: Sale{
				Account: "MIOTA",
				Amount:  10.0,
			},
			Debit: Sale{
				Account: "XRP",
				Amount:  50.0,
			},
			Date: FormattedTime{time.Now().Add(2 * time.Hour)},
		},
		{
			Credit: Sale{
				Account: "BTC",
				Amount:  1,
			},
			Debit: Sale{
				Account: "USD",
				Amount:  10000,
			},
			Date: FormattedTime{time.Now().Add(3 * time.Hour)},
		},
		{
			Credit: Sale{
				Account: "ETH",
				Amount:  100,
			},
			Debit: Sale{
				Account: "BTC",
				Amount:  0.25,
			},
			Date: FormattedTime{time.Now().Add(4 * time.Hour)},
		},
	}

	for _, rec := range recs {
		err = db.AddRecord(&rec)

		assert.Nil(t, err)
	}

	res, err := db.GetRecords()

	assert.Nil(t, err)
	assert.Equal(t, len(res), len(recs))

	count := len(recs) - 1

	for i := 0; i <= count; i++ {
		fmt.Println(res[i])
		assert.Equal(t, res[i].Debit.Account, recs[count-i].Debit.Account)
		assert.Equal(t, res[i].Debit.Amount, recs[count-i].Debit.Amount)
		assert.Equal(t, res[i].Credit.Account, recs[count-i].Credit.Account)
		assert.Equal(t, res[i].Credit.Amount, recs[count-i].Credit.Amount)
	}

	sales, err := db.AggregateRecords()

	assert.Nil(t, err)

	fmt.Println(sales)

	assert.Equal(t, 5, len(sales))

	find := func(n Unit) *Sale {
		for _, s := range sales {
			if n == s.Account {
				return &s
			}
		}
		return nil
	}

	assert.Equal(t, Sale{"BTC", 0.75}, *find("BTC"))
	assert.Equal(t, Sale{"USD", -10100}, *find("USD"))
	assert.Equal(t, Sale{"MIOTA", 10}, *find("MIOTA"))
	assert.Equal(t, Sale{"ETH", 100}, *find("ETH"))
	assert.Equal(t, Sale{"XRP", 50}, *find("XRP"))

}
