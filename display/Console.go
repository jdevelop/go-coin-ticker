// +build linux
// +build amd64 386

package display

import (
	"fmt"
)

type console struct{}

var c = console{}

func (c *console) Render(line int, text string) {
	fmt.Printf("%1d: %2s\n", line, text)
}

func (c *console) Clear() {
	//
}

func MakeDisplay(data []int, rs int, e int) (Display, error) {
	return &c, nil
}

func (c *console) PriceUp(oldPrice float64, newPrice float64) {
	fmt.Printf("⇈ %1.4f => %2.4f\n", oldPrice, newPrice)
}

func (c *console) PriceDown(oldPrice float64, newPrice float64) {
	fmt.Printf("⇊ %1.4f => %2.4f\n", oldPrice, newPrice)
}

func MakeLED(up int, down int) *console {
	return &c
}
