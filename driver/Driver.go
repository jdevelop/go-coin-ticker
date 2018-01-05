package driver

import (
	"fmt"
	"math"

	"github.com/jdevelop/go-coin-ticker/display"
	"github.com/jdevelop/go-coin-ticker/market"
	"github.com/jdevelop/go-coin-ticker/storage"
)

type history struct {
	price     float64
	timestamp int32
}

//Driver is the primary initialization structure that holds all the
//settings necessary for running the price checking and display routine.
type Driver struct {
	tickers market.TickersPipeline
	display display.Display
	db      storage.RecordsDAO
	signal  map[string]display.PriceSignal

	history map[string]history
}

//TickerUpdate method runs the fetch of the ticker data and refresh of the contents.
func (d *Driver) TickerUpdate(tickers []string) {
	for i, symbol := range tickers {
		ticker, err := d.tickers.FetchCoins(symbol)
		if err != nil {
			fmt.Println(err)
			continue
		}
		last, historyExists := d.history[ticker.Symbol]
		d.history[ticker.Symbol] = history{price: ticker.PriceUSD, timestamp: ticker.LastUpdated}
		if last.price != ticker.PriceUSD && last.timestamp < ticker.LastUpdated {
			d.display.Render(i, "                ")
			d.display.Render(i, fmt.Sprintf("%1s: $%2.2f", ticker.Symbol, ticker.PriceUSD))
			signal, ok := d.signal[ticker.Symbol]
			if ok && historyExists {
				if last.price < ticker.PriceUSD {
					signal.PriceUp(last.price, ticker.PriceUSD)
				} else {
					signal.PriceDown(last.price, ticker.PriceUSD)
				}
			}
		}
	}
}

//PortfolioUpdate runs the portfolio value update routine.
func (d *Driver) PortfolioUpdate() {

	recs, err := d.db.AggregateRecords()
	if err != nil {
		d.display.Clear()
		d.display.Render(0, "ERROR")
		return
	}

	total := 0.0
	gain := 0.0
	price := 0.0

	for _, rec := range recs {
		if storage.IsDebit(rec.Account) {
			price = price + rec.Amount
			continue
		}

		if storage.IsFee(rec.Account) {
			continue
		}

		sym, err := d.tickers.FetchCoins(rec.Account)
		if err != nil {
			continue
		}
		total = total + rec.Amount*sym.PriceUSD
		gain = gain + rec.Amount*sym.PriceUSD
	}

	price = math.Abs(price)

	d.display.Clear()
	d.display.Render(0, fmt.Sprintf("$%-5.2f/$%-5.2f", total, price))
	d.display.Render(1, fmt.Sprintf("$%+-5.2f:%-+5.1f%%", gain-price, (gain-price)*100/price))
}

//MakeDriver creates the instance of the Driver and initializes the internal fields.
func MakeDriver(
	tickers market.TickersPipeline,
	display display.Display,
	signals map[string]display.PriceSignal,
	db storage.RecordsDAO,
) *Driver {
	return &Driver{
		signal:  signals,
		history: make(map[string]history),
		display: display,
		tickers: tickers,
		db:      db,
	}
}
