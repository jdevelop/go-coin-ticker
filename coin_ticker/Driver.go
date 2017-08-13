package coin_ticker

import "fmt"

type history struct {
	price     float64
	timestamp int32
}

type Driver struct {
	tickers TickersPipeline
	display Display
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
		last, historyExists := d.history[ticker.Id]
		d.history[ticker.Id] = history{price: ticker.PriceUSD, timestamp: ticker.LastUpdated}
		if last.price != ticker.PriceUSD && last.timestamp < ticker.LastUpdated {
			d.display.Render(i, fmt.Sprintf("%1s: $%2.2f", ticker.Symbol, ticker.PriceUSD))
			signal, ok := d.signal[ticker.Id]
			if ok && historyExists {
				if last.price < ticker.PriceUSD {
					signal.priceUp(ticker.PriceUSD)
				} else {
					signal.priceDown(ticker.PriceUSD)
				}
			}
		}
	}
}

func MakeDriver(tickers TickersPipeline, display Display, signals map[string]PriceSignal) *Driver {
	dr := Driver{
		signal:  signals,
		history: make(map[string]history),
		display: display,
		tickers: tickers,
	}
	return &dr
}
