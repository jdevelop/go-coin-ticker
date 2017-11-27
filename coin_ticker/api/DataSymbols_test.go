package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"time"
)

const TmpPath = "/tmp/coins.db.test"

func Test_DBAccess(t *testing.T) {

	os.Remove(TmpPath)

	db, err := MakeDB(TmpPath)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	recs := []Record{
		{
			Price:  1.00055,
			Amount: 0.1128,
			Symbol: "ETH",
			Date:   time.Now(),
		},
		{
			Price:  8000.221,
			Amount: 1,
			Symbol: "BTC",
			Date:   time.Now().Add(1 * time.Hour),
		},
		{
			Price:  15,
			Amount: 1,
			Symbol: "ETH",
			Date:   time.Now().Add(2 * time.Hour),
		},
		{
			Price:  9000,
			Amount: 17,
			Symbol: "BTC",
			Date:   time.Now().Add(3 * time.Hour),
		},
		{
			Price:  1,
			Amount: 1,
			Symbol: "XRP",
			Date:   time.Now().Add(4 * time.Hour),
		},
	}

	for _, rec := range recs {
		err = db.AddRecord(&rec)

		assert.Nil(t, err)
	}

	res, err := db.GetRecords()

	assert.Nil(t, err)

	for i := 0; i <= 4; i++ {
		assert.Equal(t, res[i].Symbol, recs[4-i].Symbol)
		assert.Equal(t, res[i].Amount, recs[4-i].Amount)
		assert.Equal(t, res[i].Price, recs[4-i].Price)
	}

	res, err = db.AggregateRecords()

	assert.Nil(t, err)

	assert.Equal(t, 3, len(res))

}
