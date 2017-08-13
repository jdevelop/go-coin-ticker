package coin_ticker

import (
	"fmt"
	"github.com/davecheney/gpio"
)

type PriceSignal interface {
	PriceUp(oldPrice float64, newPrice float64)
	PriceDown(oldPrice float64, newPrice float64)
	Clear()
}

func (c *console) priceUp(oldPrice float64, newPrice float64) {
	fmt.Printf("⇈ %1.4f => %2.4f\n", oldPrice, newPrice)
}

func (c *console) priceDown(oldPrice float64, newPrice float64) {
	fmt.Printf("⇊ %1.4f => %2.4f\n", oldPrice, newPrice)
}

func (c *console) Clear() {}

type LED struct {
	pinUp   gpio.Pin
	pinDown gpio.Pin
}

func MakeLED(up int, down int) *LED {
	pinUp, _ := gpio.OpenPin(up, gpio.ModeOutput)
	pinDown, _ := gpio.OpenPin(down, gpio.ModeOutput)
	return &LED{
		pinUp:   pinUp,
		pinDown: pinDown,
	}
}

func (l *LED) PriceUp(oldPrice float64, newPrice float64) {
	l.pinDown.Clear()
	l.pinUp.Set()
}

func (l *LED) PriceDown(oldPrice float64, newPrice float64) {
	l.pinDown.Set()
	l.pinUp.Clear()
}

func (l *LED) Clear() {
	l.pinDown.Clear()
	l.pinUp.Clear()
}
