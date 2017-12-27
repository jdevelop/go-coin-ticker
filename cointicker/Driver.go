package cointicker

import (
	"fmt"
	"math"
)

type history struct {
	price     float64
	timestamp int32
}

type Driver struct {
	tickers TickersPipeline
	display Display
	db      RecordsDAO
	signal  map[string]PriceSignal

	history map[string]history
}

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
		if IsDebit(rec.Account) {
			price = price + rec.Amount
			continue
		}

		if IsFee(rec.Account) {
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

func MakeDriver(
	tickers TickersPipeline,
	display Display,
	signals map[string]PriceSignal,
	db RecordsDAO,
) *Driver {
	return &Driver{
		signal:  signals,
		history: make(map[string]history),
		display: display,
		tickers: tickers,
		db:      db,
	}
}
