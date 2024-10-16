package main

import (
	"os"
	"saturn/assembler"
	"saturn/gui"
	"saturn/linker"
)

func main() {
	programs := os.Args[1:]
	if len(programs) == 0 { // default
		programs = append(programs, "saturn_test1.asm")
		programs = append(programs, "saturn_test2.asm")
	}

	stackLimit, programName := linker.Run(assembler.Run(programs...))
	gui.Initialize(stackLimit)
	gui.LoadProgram(ReadProgram(programName + ".hpx"))
	gui.Run()
}
