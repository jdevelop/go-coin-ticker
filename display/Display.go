package display

//Display interface for rendering some data either in console or using LCD or other implementation.
type Display interface {
	Render(line int, text string)
	Clear()
}
