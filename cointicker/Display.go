package cointicker

import (
	"fmt"

	lcd "github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
)

type Display interface {
	Render(line int, text string)
	Clear()
}

type console struct{}

func (c *console) Render(line int, text string) {
	fmt.Printf("%1d: %2s\n", line, text)
}

func MakeConsoleDisplay() Display {
	return &console{}
}

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

func MakeLCDDisplay(data []int, rs int, e int) (d Display, err error) {
	lcdRef, err := lcd.NewLCD4(data, rs, e)
	if err == nil {
		lcdRef.Init()
		d = &LCD{lcdRef: &lcdRef}
	}
	return
}
