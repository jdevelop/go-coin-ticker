package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jdevelop/go-coin-ticker/coin_ticker"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"flag"
)

type Config struct {
	Ticker struct {
		LCD struct {
			DataPins []int `mapstructure:"data-pins"`
			RsPin    int   `mapstructure:"rs-pin"`
			EPin     int   `mapstructure:"e-pin"`
		} `mapstructure:"lcd"`
		DB struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"db"`
		Symbols []string `mapstructure:"symbols"`
		REST struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"rest"`
	} `mapstructure:"coin-ticker"`
}

func main() {

	portfolio := flag.Bool("portfolio", true, "")

	flag.Parse()

	var display coin_ticker.Display
	signals := make(map[string]coin_ticker.PriceSignal)

	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath("$HOME/.coins_ticker") // call multiple times to add many search paths
	err := viper.ReadInConfig()                // Find and read the config file
	conf := Config{}
	err = viper.Unmarshal(&conf)
	if err != nil || len(conf.Ticker.LCD.DataPins) == 0 {
		display = coin_ticker.MakeConsoleDisplay()
	} else {
		display, err = coin_ticker.MakeLCDDisplay(
			conf.Ticker.LCD.DataPins,
			conf.Ticker.LCD.RsPin,
			conf.Ticker.LCD.EPin,
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

	market := coin_ticker.MakeCoinMarket()

	db, err := coin_ticker.MakeDB(conf.Ticker.DB.Path)

	driver := coin_ticker.MakeDriver(
		market,
		display,
		signals,
		db,
	)

	if conf.Ticker.REST.Host != "" {

		if err != nil {
			log.Fatal(err)
		}

		r := coin_ticker.MakeREST(db, market)
		go func() {
			addr := fmt.Sprintf("%s:%d", conf.Ticker.REST.Host, conf.Ticker.REST.Port)
			fmt.Printf("Starting REST at %s\n", addr)
			http.ListenAndServe(addr, r)
		}()
	}

	var upd func()

	if *portfolio {
		upd = func() { driver.PortfolioUpdate() }
	} else {
		upd = func() { driver.TickerUpdate(conf.Ticker.Symbols) }
	}
	ticker := time.Tick(10 * time.Second)
	upd()

	fmt.Printf("Starting ticker on %v\n", conf.Ticker.Symbols)

	for range ticker {
		upd()
	}
}
