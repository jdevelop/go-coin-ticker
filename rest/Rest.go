package rest

import (
	"encoding/json"
	"github.com/jdevelop/go-coin-ticker/market"
	"github.com/jdevelop/go-coin-ticker/storage"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mgutz/logxi/v1"
)

//PricedRecord contains the market symbol with price and quantity.
type PricedRecord struct {
	Symbol      storage.Unit `json:"symbol,omitempty"`
	MarketPrice float64      `json:"market_price,omitempty"`
	Qty         float64      `json:"qty,omitempty"`
	Value       float64      `json:"value,omitempty"`
}

//Dashboard holds the transaction data aggregated by coin symbol.
type Dashboard struct {
	TotalReturn float64        `json:"total_return"`
	TotalSpent  float64        `json:"total_spent"`
	GainLoss    float64        `json:"gain_loss"`
	Symbols     []PricedRecord `json:"symbols"`
}

func jsonCT(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func httpError(resp http.ResponseWriter, msg string, code int, err error) {
	log.Error(msg, err)
	http.Error(resp, msg, code)
}

//MakeREST creates the REST interface with provided backends for the data storage and market access.
func MakeREST(db storage.RecordsDAO, m market.TickersPipeline) (r *httprouter.Router) {
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

		codes := make([]string, 0)
		for _, r := range res {
			codes = append(codes, r.Account)
		}

		symMap, err := m.FetchCoins(codes...)
		if err != nil {
			httpError(w, "Can't process coin", 500, err)
			return
		}

		for _, rec := range res {
			if storage.IsDebit(rec.Account) {
				spent = spent + rec.Amount
				continue
			}
			if storage.IsFee(rec.Account) {
				continue
			}
			sym, ok := symMap[rec.Account]
			if !ok {
				log.Error("Can't read market price for ", sym)
				continue
			}
			total = total + sym.PriceUSD*rec.Amount
			recs = append(recs, PricedRecord{
				Symbol:      rec.Account,
				Qty:         rec.Amount,
				MarketPrice: sym.PriceUSD,
				Value:       sym.PriceUSD * rec.Amount,
			})
		}

		spent = math.Abs(spent)

		gain := 0.0
		if spent != 0 {
			gain = (total - spent) * 100 / spent
		}
		data, err := json.Marshal(&Dashboard{
			TotalReturn: total,
			TotalSpent:  spent,
			GainLoss:    gain,
			Symbols:     recs,
		})

		if err != nil {
			httpError(w, "Can't serialize JSON", 500, err)
			return
		}
		jsonCT(w)
		w.Write(data)
	})

	r.GET("/list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res, err := db.GetRecords()
		if err != nil {
			httpError(w, "Can't load data from DB", 500, err)
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			httpError(w, "No data", 500, err)
			return
		}
		jsonCT(w)
		w.Write(data)
	})

	r.DELETE("/remove/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idStr := p.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpError(w, "Not an ID", 400, err)
			return
		}
		if err = db.RemoveRecord(id); err != nil {
			httpError(w, "Can't remove record", 500, err)
			return
		}
		jsonCT(w)
		w.Write([]byte("{ 'status' : 'Complete' }"))
	})

	r.PUT("/transfer", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpError(w, "Unable to read request", 400, err)
			return
		}
		var rec storage.Record
		if err = json.Unmarshal(data, &rec); err != nil {
			log.Error("Wrong record", string(data))
			httpError(w, "Unable to unmarshal request", 400, err)
			return
		}

		debitAct := strings.ToLower(rec.Debit.Account)
		creditAct := strings.ToLower(rec.Credit.Account)

		coins, err := m.FetchCoins(debitAct, creditAct)
		if err != nil {
			httpError(w, "Can't communicate with market", 400, err)
			return
		}

		if (!storage.IsDebit(debitAct) && coins[debitAct].ID == "") ||
			(!storage.IsDebit(creditAct) && coins[creditAct].ID == "") {
			httpError(w, "coin not found", 400, err)
			return
		}

		if rec.Credit.Amount <= 0 || rec.Debit.Amount <= 0 {
			httpError(w, "amount must be positive", 400, err)
			return
		}

		if err = db.AddRecord(&rec); err != nil {
			httpError(w, "Unable to save record", 400, err)
			return
		}
		jsonCT(w)
		w.Write([]byte("{ status: 'Complete' }"))
	})

	return
}
