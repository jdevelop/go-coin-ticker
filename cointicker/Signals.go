package cointicker

type PriceSignal interface {
	PriceUp(oldPrice float64, newPrice float64)
	PriceDown(oldPrice float64, newPrice float64)
	Clear()
}
