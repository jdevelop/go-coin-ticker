package display

type Display interface {
	Render(line int, text string)
	Clear()
}
