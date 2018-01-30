// +build linux,arm

package display

import (
	"github.com/davecheney/gpio"
	lcd "github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
)

//LCD defines the properties to be used for LCD screen (pinout)
type LCD struct {
	lcdRef lcd.PiLCD
}

func (lcd *LCD) Render(line int, text string) {
	lcd.lcdRef.SetCursor(uint8(line), 0)
	lcd.lcdRef.Print(text)
}

func (lcd *LCD) Clear() {
	lcd.lcdRef.Cls()
}

func MakeDisplay(data []int, rs int, e int) (d Display, err error) {
	lcdRef, err := lcd.NewLCD4(data, rs, e)
	if err == nil {
		lcdRef.Init()
		d = &LCD{lcdRef: &lcdRef}
	}
	return
}

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
