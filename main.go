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
		programs = append(programs, "linker/linker_test.asm")
		programs = append(programs, "linker/linker_test_part2.asm")
	}

	stackLimit, programName := linker.Run(assembler.Run(programs...))
	gui.Initialize(stackLimit)
	gui.LoadProgram(ReadProgram(programName + ".hpx"))
	gui.Run()
}
