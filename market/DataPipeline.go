package market

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	ticker "github.com/jdevelop/go-coin-ticker/cointicker"
)

type TickersPipeline interface {
	FetchCoins(symbol string) (*ticker.TickerData, error)
}

type coinMarket struct{}

func MakeCoinMarket() TickersPipeline {
	return &coinMarket{}
}

func (mkt *coinMarket) FetchCoins(coinCode string) (*ticker.TickerData, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + coinCode + "/")
	if err != nil {
		ticker.LogError(err)
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ticker.LogError(err)
		return nil, err
	}
	var tickers []ticker.TickerData
	err = json.Unmarshal(res, &tickers)
	if err != nil {
		ticker.LogError(err)
		return nil, err
	}
	return &tickers[0], nil
}
