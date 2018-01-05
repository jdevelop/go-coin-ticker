package market

type TickerData struct {
	Id               string  `json:"id"`
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
