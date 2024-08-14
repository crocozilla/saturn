package gui

import (
	"fmt"
	"saturn/shared"
	"saturn/vm"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var program shared.Program = shared.Program{
	shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.ADD,
		Operands:    shared.Operands{First: 30},
	},
}

func Run(machine *vm.VirtualMachine) {
	a := app.New()

	left := container.NewVBox(buttons(machine))
	middle := container.NewVBox(registers(machine))
	right := container.NewVBox(memory(machine))

	root := container.NewHBox(left, layout.NewSpacer(), middle, layout.NewSpacer(), right)

	w := a.NewWindow("Saturn")
	w.Resize(fyne.NewSize(1200, 700))
	w.SetContent(root)

	w.ShowAndRun()
}

// func update(d *fyne.Container) {
	
// }

func registers(vm *vm.VirtualMachine) *fyne.Container {
	r := container.NewGridWithColumns(3)

	r.Add(widget.NewLabel(fmt.Sprintf("Program Counter: %d", vm.PC())))
	r.Add(widget.NewLabel(fmt.Sprintf("Stack Pointer: %d", vm.SP())))
	r.Add(widget.NewLabel(fmt.Sprintf("Acumulador: %d", vm.Acumulator())))
	r.Add(widget.NewLabel(fmt.Sprintf("Operação: %d", vm.Operation())))
	r.Add(widget.NewLabel(fmt.Sprintf("Endereço de Memória: %d", vm.MemoryAddress())))

	go func() {
		for range time.Tick(time.Millisecond) {
			r.Refresh()
		}
	}()

	return r
}
