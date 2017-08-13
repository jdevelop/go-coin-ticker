package coin_ticker

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var dataCoin = `{ "id": "ethereum",
"name": "Ethereum",
"symbol": "ETH",
"rank": "2",
"price_usd": "308.881",
"price_btc": "0.084432",
"24h_volume_usd": "938420000.0",
"market_cap_usd": "29014763385.0",
"available_supply": "93935086.0",
"total_supply": "93935086.0",
"percent_change_1h": "0.37",
"percent_change_24h": "3.79",
"percent_change_7d": "37.82",
"last_updated": "1502496849"
}`

var testTicker = TickerData{
	Id:               "ethereum",
	Name:             "Ethereum",
	Symbol:           "ETH",
	Rank:             2,
	PriceUSD:         308.881,
	PriceBTC:         0.084432,
	Volume24H:        938420000.0,
	MarketCapUSD:     29014763385.0,
	AvailableSupply:  93935086.0,
	TotalSupply:      93935086.0,
	PercentChange1H:  0.37,
	PercentChange24H: 3.79,
	PercentChange7D:  37.82,
	LastUpdated:      1502496849,
}

func TestJsonUnmarshalling(t *testing.T) {
	z := &TickerData{}

	json.Unmarshal([]byte(dataCoin), z)

	assert.Equal(t, testTicker, *z)
}

func TestJsonUnmarshallingArray(t *testing.T) {

	var coins []TickerData

	json.Unmarshal([]byte("["+dataCoin+"]"), &coins)

	assert.Equal(t, testTicker, coins[0])
}
