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

var program shared.Program = shared.Program{
	shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.ADD,
		Operands:    shared.Operands{First: 30, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.DIRECT_IMMEDIATE,
		Operation:   shared.COPY,
		Operands:    shared.Operands{First: 64, Second: 5},
	},
	shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.INJ,
		Operands:    shared.Operands{First: 64, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.INDIRECT,
		Operation:   shared.DIVIDE,
		Operands:    shared.Operands{First: 0, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.ADD,
		Operands:    shared.Operands{First: 64, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.STORE,
		Operands:    shared.Operands{First: 72, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.INDIRECT,
		Operation:   shared.LOAD,
		Operands:    shared.Operands{First: 0, Second: 0},
	},
	shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.STOP,
		Operands:    shared.Operands{},
	},
}

var machine = vm.New()

var mem = container.NewGridWithColumns(4)
var r = container.NewGridWithColumns(3)

func InsertProgram(program []shared.Word) {
	machine.InsertProgram(program)
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
		machine.ExecuteAll(program)
		updateGUI()
	})

	resetBtn := widget.NewButton("Resetar", func() {
		machine.Reset()
		updateGUI()
	})

	return container.New(layout.NewGridLayout(3), executeBtn, executeAllBtn, resetBtn)
}
