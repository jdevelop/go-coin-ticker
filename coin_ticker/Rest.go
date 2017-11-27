package coin_ticker

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/mgutz/logxi/v1"
	"encoding/json"
	"io/ioutil"
	"time"
	"strconv"
)

type PricedRecord struct {
	Record
	MarketPrice float64 `json:"market_price"`
	Return      float64 `json:"return"`
	Value       float64 `json:"value"`
}

type Dashboard struct {
	TotalReturn float64        `json:"total_return"`
	TotalSpent  float64        `json:"total_spent"`
	GainLoss    float64        `json:"gain_loss"`
	Symbols     []PricedRecord `json:"symbols"`
}

func jsonCT(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func MakeREST(db RecordsDAO, m TickersPipeline) (r *httprouter.Router) {
	r = httprouter.New()

	r.GET("/dashboard", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res, err := db.AggregateRecords()
		if err != nil {
			http.Error(w, "No data", 400)
			log.Error("Can't read the database", err)
			return
		}
		recs := make([]PricedRecord, 0)
		total := 0.0
		spent := 0.0
		for _, rec := range res {
			sym, err := m.FetchCoins(rec.Symbol)
			if err != nil {
				log.Error("Can't process coin", err)
				continue
			}
			total = total + sym.PriceUSD*rec.Amount - rec.Price
			spent = spent + rec.Price
			recs = append(recs, PricedRecord{
				Record:      rec,
				MarketPrice: sym.PriceUSD,
				Value:       sym.PriceUSD * rec.Amount,
				Return:      sym.PriceUSD*rec.Amount - rec.Price,
			})
		}

		data, err := json.Marshal(&Dashboard{
			TotalReturn: total,
			TotalSpent:  spent,
			GainLoss:    total / spent * 100,
			Symbols:     recs,
		})

		if err != nil {
			http.Error(w, "Unknown format", 500)
			return
		}
		jsonCT(w)
		w.Write(data)
	})

	r.GET("/list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res, err := db.GetRecords()
		if err != nil {
			log.Error("Can't read data", err)
			http.Error(w, "Can't load data from DB", 500)
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			log.Error("Can't marshal JSON data", err)
			http.Error(w, "No data", 500)
			return
		}
		jsonCT(w)
		w.Write(data)
	})

	r.DELETE("/remove/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idStr := p.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("Can't convert ID", err)
			http.Error(w, "Not an ID", 400)
			return
		}
		if err = db.RemoveRecord(id); err != nil {
			log.Error("Can't remove ID", id, err)
			http.Error(w, "Can't remove record", 500)
			return
		}
		jsonCT(w)
		w.Write([]byte("{ 'status' : 'Complete' }"))
	})

	r.PUT("/buy", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request", err)
			http.Error(w, "Unable to read request", 400)
			return
		}
		var rec Record
		if err = json.Unmarshal(data, &rec); err != nil {
			log.Error("Can't unmarshal request", err)
			http.Error(w, "Unable to unmarshal request", 400)
			return
		}
		rec.Date = time.Now()
		if err = db.AddRecord(&rec); err != nil {
			log.Error("Can't save record", err)
			http.Error(w, "Unable to save record", 400)
			return
		}
		jsonCT(w)
		w.Write([]byte("{ status: 'Complete' }"))
	})

	return
}
