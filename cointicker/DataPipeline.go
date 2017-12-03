package cointicker

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

func (mkt *coinMarket) FetchCoins(coinCode string) (*TickerData, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + coinCode + "/")
	if err != nil {
		LogError(err)
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError(err)
		return nil, err
	}
	var tickers []TickerData
	err = json.Unmarshal(res, &tickers)
	if err != nil {
		LogError(err)
		return nil, err
	}
	return &tickers[0], nil
}
