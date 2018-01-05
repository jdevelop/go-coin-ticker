package market

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

//TickerData holds the JSON-compatible object of the current coin data from coinmarketcap.com
type TickerData struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int8    `json:"rank,string"`
	PriceUSD         float64 `json:"price_usd,string"`
	PriceBTC         float64 `json:"price_btc,string"`
	Volume24H        float64 `json:"24h_volume_usd,string"`
	MarketCapUSD     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1H  float32 `json:"percent_change_1h,string"`
	PercentChange24H float32 `json:"percent_change_24h,string"`
	PercentChange7D  float32 `json:"percent_change_7d,string"`
	LastUpdated      int32   `json:"last_updated,string"`
}

//TickersPipeline defines the methods to fetch the coin stats data.
type TickersPipeline interface {
	FetchCoins(symbol ...string) (map[string]TickerData, error)
}

type coinMarket struct{}

//MakeCoinMarket creates the default instance of the coin market interface.
func MakeCoinMarket() TickersPipeline {
	return &coinMarket{}
}

func (mkt *coinMarket) FetchCoins(coinCode ...string) (t map[string]TickerData, err error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
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
	t = make(map[string]TickerData)
	for _, c := range coinCode {
		t[strings.ToLower(c)] = TickerData{}
	}
	for _, tck := range tickers {
		n := strings.ToLower(tck.Symbol)
		if _, ok := t[n]; ok {
			t[n] = tck
		}
	}
	return
}
