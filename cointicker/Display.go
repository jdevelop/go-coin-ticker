package cointicker

type Display interface {
	Render(line int, text string)
	Clear()
}
