package coin_ticker

import (
	"fmt"
	"github.com/davecheney/gpio"
)

type PriceSignal interface {
	priceUp(newPrice float64)
	priceDown(newPrice float64)
}

func (c *console) priceUp(newPrice float64) {
	fmt.Printf("⇈ %1.4f\n", newPrice)
}

func (c *console) priceDown(newPrice float64) {
	fmt.Printf("⇊ %1.4f\n", newPrice)
}

type LED struct {
	pinUp   gpio.Pin
	pinDown gpio.Pin
}

func (l *LED) priceUp(newPrice float64) {
	l.pinDown.Clear()
	l.pinUp.Set()
}

func (l *LED) priceDown(newPrice float64) {
	l.pinDown.Set()
	l.pinUp.Clear()
}
