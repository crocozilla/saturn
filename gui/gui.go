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

var mem = container.NewGridWithColumns(4)
var r = container.NewGridWithColumns(3)
var machine *vm.VirtualMachine
var output *widget.Label
var programBackup []shared.Word

func Initialize(stackLimit uint16) {
	machine = vm.New(stackLimit)
	output = widget.NewLabel(strconv.Itoa((int(machine.Output()))))
}

func LoadProgram(program []shared.Word) {
	programBackup = program
	machine.LoadProgram(program)
}

func ReInsertProgram() {
	machine.LoadProgram(programBackup)
}

func Run() {
	a := app.New()

	left := container.NewVBox(buttons())
	middle := container.NewVBox(registers(), io(), buttons())
	right := container.NewVBox(memory())

	root := container.NewHBox(left, layout.NewSpacer(), middle, layout.NewSpacer(), right)

	w := a.NewWindow("Saturn")
	w.Resize(fyne.NewSize(1200, 700))
	w.SetContent(root)

	updateGUI()

	w.ShowAndRun()
}

func registers() fyne.Widget {
	return widget.NewCard("Registradores", "", r)
}

func updateGUI() {
	output.SetText(strconv.Itoa((int(machine.Output()))))

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

func memory() fyne.Widget {
	scrollable := container.NewVScroll(mem)
	scrollable.SetMinSize(fyne.NewSize(300, 700))

	backgroundColor := color.RGBA{R: 0, B: 0, G: 0, A: 50}
	background := canvas.NewRectangle(backgroundColor)

	withBackground := container.NewStack(background, scrollable)
	return widget.NewCard("Memória", "", withBackground)
}

func buttons() *fyne.Container {
	executeBtn := widget.NewButton("Executar", func() {
		if machine.IsRunning() {
			machine.Execute()
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

	return container.NewVBox(executeBtn, container.NewHBox(executeAllBtn, resetBtn))
}

func io() *fyne.Container {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Digite a entrada")
	inputBtn := widget.NewButton("Salvar", func() {
		data, err := strconv.Atoi(inputEntry.Text)
		if err != nil {
			panic(err)
		}
		machine.SetInput(uint16(data))
	})

	input := widget.NewCard("Entrada", "", container.NewGridWithColumns(2, inputEntry, inputBtn))
	output := widget.NewCard("Saída", "", container.NewStack(canvas.NewRectangle(color.Black), output))

	return container.NewVBox(input, output)
}
