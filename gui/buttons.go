package gui

import (
	"fmt"
	"saturn/shared"
	"saturn/vm"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buttons(machine *vm.VirtualMachine) *fyne.Container {
	executeAllBtn := widget.NewButton("Executar", func() {
		fmt.Println("Executar tudo")
	})

	executeBtn := widget.NewButton("Executar Tudo", func() {
		machine.Execute(shared.Instruction{
			Operation:   shared.ADD,
			Operands:    shared.Operands{First: 10, Second: 0},
			AddressMode: shared.IMMEDIATE,
		})

		fmt.Println(machine.Acumulator())
	})

	return container.New(layout.NewGridLayout(3), executeBtn, executeAllBtn)
}
