package main

import (
	"saturn/gui"
	"saturn/vm"
)

func main() {
	machine := vm.New()

	gui.Run(machine)
}
