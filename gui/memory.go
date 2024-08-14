package gui

import (
	"image/color"
	"saturn/vm"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func memory(machine *vm.VirtualMachine) fyne.CanvasObject {
	mem := container.NewGridWithColumns(4)


	for i, value := range machine.Memory() {
		textAddress := canvas.NewText(strconv.Itoa(int(i)), color.White)
		textValue := canvas.NewText("[" + strconv.Itoa(int(value)) + "]", color.RGBA{R: 255, B: 0, G: 255, A: 255})
		cont := container.NewHBox(textAddress, layout.NewSpacer(), textValue)

		mem.Add(cont)
	}

	scrollable := container.NewVScroll(mem)
	scrollable.SetMinSize(fyne.NewSize(300, 700))

	backgroundColor := color.RGBA{R: 0, B: 0, G: 0, A: 50}
	background := canvas.NewRectangle(backgroundColor)

	withBackground := container.NewStack(background, scrollable)
	return withBackground
}
