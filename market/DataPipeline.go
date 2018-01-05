package market

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TickersPipeline interface {
	FetchCoins(symbol string) (*TickerData, error)
}

type coinMarket struct{}

func MakeCoinMarket() TickersPipeline {
	return &coinMarket{}
}

func (mkt *coinMarket) FetchCoins(coinCode string) (t *TickerData, err error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + coinCode + "/")
	if err != nil {
		return
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var tickers []TickerData
	err = json.Unmarshal(res, &tickers)
	if err != nil {
		return
	}
	t = &tickers[0]
	return
}
