package main

import (
	"github.com/jdevelop/go-coin-ticker/coin_ticker"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
	"fmt"
)

func main() {

	var display coin_ticker.Display
	signals := make(map[string]coin_ticker.PriceSignal)

	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath("$HOME/.coins_ticker") // call multiple times to add many search paths
	err := viper.ReadInConfig()                // Find and read the config file
	if err != nil || viper.GetString("coin-ticker.lcd.data-pins") == "" {
		display = coin_ticker.MakeConsoleDisplay()
	} else {
		dataPinsStr := strings.Split(viper.GetString("coin-ticker.lcd.data-pins"), ",")
		dataPins := make([]int, len(dataPinsStr))
		for i, v := range dataPinsStr {
			dataPins[i], _ = strconv.Atoi(v)
		}
		display, err = coin_ticker.MakeLCDDisplay(
			dataPins,
			viper.GetInt("coin-ticker.lcd.rs-pin"),
			viper.GetInt("coin-ticker.lcd.e-pin"),
		)
		asInt := func(s interface{}) int {
			switch s.(type) {
			case string:
				res, _ := strconv.Atoi(s.(string))
				return res
			case int:
				return s.(int)
			default:
				return -1
			}
		}
		for coins, v := range viper.GetStringMap("coin-ticker.pins") {
			coinData := v.(map[string]interface{})
			signals[strings.ToUpper(coins)] = coin_ticker.MakeLED(asInt(coinData["pin-up"]), asInt(coinData["pin-down"]))
		}

		display.Render(0, "  COINS  ")
		display.Render(1, " TRACKER ")

		delay := 2 * time.Second

		for k, signal := range signals {
			display.Render(0, fmt.Sprintf("Testing %1s UP", k))
			signal.PriceUp(0, 0)
			time.Sleep(delay)
			display.Render(1, fmt.Sprintf("Testing %1s DOWN", k))
			signal.PriceDown(0, 0)
			time.Sleep(delay)
			signal.Clear()
		}

		display.Clear()

	}

	driver := coin_ticker.MakeDriver(
		coin_ticker.MakeCoinMarket(),
		display,
		signals,
	)

	tickers := []string{"ethereum", "bitcoin"}

	driver.TickerUpdate(tickers)
	ticker := time.Tick(10 * time.Second)
	for range ticker {
		driver.TickerUpdate(tickers)
	}
}
