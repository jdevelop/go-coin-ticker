package main

import (
	"fmt"
	"github.com/jdevelop/go-coin-ticker/driver"
	"strconv"
	"strings"
	"time"

	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/jdevelop/go-coin-ticker/display"
	"github.com/jdevelop/go-coin-ticker/market"
	"github.com/jdevelop/go-coin-ticker/rest"
	"github.com/jdevelop/go-coin-ticker/storage"
	"github.com/spf13/viper"
)

// Config defines the configuration file structure.
type Config struct {
	Ticker struct {
		Interval time.Duration `mapstructure:"tick"`
		LCD      struct {
			DataPins []int `mapstructure:"data-pins"`
			RsPin    int   `mapstructure:"rs-pin"`
			EPin     int   `mapstructure:"e-pin"`
		} `mapstructure:"lcd"`
		DB struct {
			Path string `mapstructure:"path"`
		} `mapstructure:"db"`
		Symbols []string `mapstructure:"symbols"`
		REST    struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"rest"`
	} `mapstructure:"coin-ticker"`
}

func main() {

	portfolio := flag.Bool("portfolio", true, "")
	noTicker := flag.Bool("noticker", false, "")

	flag.Parse()

	var dsp display.Display

	signals := make(map[string]display.PriceSignal)

	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath("$HOME/.coins_ticker") // call multiple times to add many search paths
	err := viper.ReadInConfig()                // Find and read the config file
	conf := Config{}
	err = viper.Unmarshal(&conf)
	if err != nil || len(conf.Ticker.LCD.DataPins) == 0 {
		dsp, err = display.MakeDisplay(nil, -1, -1)
	} else {
		fmt.Println("Starting with LCD: ", conf.Ticker.LCD.DataPins, "RS:",
			conf.Ticker.LCD.RsPin, "E:", conf.Ticker.LCD.EPin)
		dsp, err = display.MakeDisplay(
			conf.Ticker.LCD.DataPins,
			conf.Ticker.LCD.RsPin,
			conf.Ticker.LCD.EPin,
		)
		if err != nil {
			log.Fatal(err)
		}
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
			signals[strings.ToUpper(coins)] = display.MakeLED(asInt(coinData["pin-up"]), asInt(coinData["pin-down"]))
		}

		dsp.Render(0, "  COINS  ")
		dsp.Render(1, " TRACKER ")

		delay := 2 * time.Second

		for k, signal := range signals {
			dsp.Render(0, fmt.Sprintf("Testing %1s UP", k))
			signal.PriceUp(0, 0)
			time.Sleep(delay)
			dsp.Render(1, fmt.Sprintf("Testing %1s DOWN", k))
			signal.PriceDown(0, 0)
			time.Sleep(delay)
			signal.Clear()
		}

		dsp.Clear()

	}

	market := market.MakeCoinMarket()

	db, err := storage.MakeDB(conf.Ticker.DB.Path)

	drv := driver.MakeDriver(
		market,
		dsp,
		signals,
		db,
	)

	if conf.Ticker.REST.Host != "" {

		if err != nil {
			log.Fatal(err)
		}

		r := rest.MakeREST(db, market)
		go func() {
			addr := fmt.Sprintf("%s:%d", conf.Ticker.REST.Host, conf.Ticker.REST.Port)
			fmt.Printf("Starting REST at %s\n", addr)
			http.ListenAndServe(addr, r)
		}()
	}

	var upd func()

	if *portfolio {
		fmt.Println("Portfolio")
		upd = func() { drv.PortfolioUpdate() }
	} else {
		fmt.Println("Ticker")
		upd = func() { drv.TickerUpdate(conf.Ticker.Symbols) }
	}

	if conf.Ticker.Interval == 0 {
		conf.Ticker.Interval = 10 * time.Second
	}

	if !*noTicker {
		ticker := time.Tick(conf.Ticker.Interval)
		upd()

		fmt.Println("Starting ticker")

		for range ticker {
			upd()
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	}
}
