package main

import (
	"github.com/jdevelop/go-coin-ticker/coin_ticker"
)

func main() {
	market := coin_ticker.MakeCoinMarket()
	driver := coin_ticker.MakeDriver(
		&market,
		coin_ticker.MakeConsoleDisplay(),
		make(map[string]coin_ticker.PriceSignal),
	)

	driver.TickerUpdate([]string{"ethereum", "bitcoin"})
}
