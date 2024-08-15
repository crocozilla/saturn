package gui

import (
	"fmt"
	"image/color"
	"saturn/shared"
	"saturn/vm"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var machine = vm.New()

var mem = container.NewGridWithColumns(4)
var r = container.NewGridWithColumns(3)
var program_backup []shared.Word

func InsertProgram(program []shared.Word) {
	program_backup = program
	machine.InsertProgram(program)
}

func ReInsertProgram() {
	machine.InsertProgram(program_backup)
}

func Run() {
	a := app.New()

	left := container.NewVBox(buttons())
	middle := container.NewVBox(r)
	right := container.NewVBox(memory())

	root := container.NewHBox(left, layout.NewSpacer(), middle, layout.NewSpacer(), right)

	w := a.NewWindow("Saturn")
	w.Resize(fyne.NewSize(1200, 700))
	w.SetContent(root)

	updateGUI()

	w.ShowAndRun()
}

func updateGUI() {
	r.RemoveAll()
	r.Add(widget.NewLabel(fmt.Sprintf("Program Counter: %d", machine.PC())))
	r.Add(widget.NewLabel(fmt.Sprintf("Stack Pointer: %d", machine.SP())))
	r.Add(widget.NewLabel(fmt.Sprintf("Acumulador: %d", machine.Accumulator())))
	r.Add(widget.NewLabel(fmt.Sprintf("Operação: %d", machine.Operation())))
	r.Add(widget.NewLabel(fmt.Sprintf("Endereço de Memória: %d", machine.MemoryAddress())))

	mem.RemoveAll()
	for i, value := range machine.Memory() {
		textAddress := canvas.NewText(strconv.Itoa(int(i)), color.White)
		textValue := canvas.NewText("["+strconv.Itoa(int(value))+"]", color.RGBA{R: 255, B: 0, G: 255, A: 255})
		cont := container.NewHBox(textAddress, layout.NewSpacer(), textValue)

		mem.Add(cont)
	}
}

func memory() fyne.CanvasObject {
	scrollable := container.NewVScroll(mem)
	scrollable.SetMinSize(fyne.NewSize(300, 700))

	backgroundColor := color.RGBA{R: 0, B: 0, G: 0, A: 50}
	background := canvas.NewRectangle(backgroundColor)

	withBackground := container.NewStack(background, scrollable)
	return withBackground
}

func buttons() *fyne.Container {
	executeBtn := widget.NewButton("Executar", func() {
		if machine.IsRunning() {

			machine.Execute(machine.PC())
			updateGUI()
		}
	})

	executeAllBtn := widget.NewButton("Executar Tudo", func() {
		machine.ExecuteAll()
		updateGUI()
	})

	resetBtn := widget.NewButton("Resetar", func() {
		machine.Reset()
		updateGUI()
	})

	return container.New(layout.NewGridLayout(3), executeBtn, executeAllBtn, resetBtn)
}
