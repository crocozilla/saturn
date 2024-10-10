package main

import (
	"saturn/assembler"
	"saturn/gui"
	"saturn/linker"
)

func main() {
	stackLimit, programName := linker.Run(assembler.Run(
		"linker/linker_test.asm", "linker/linker_test_part2.asm"))
	gui.Initialize(stackLimit)
	gui.LoadProgram(ReadProgram(programName + ".hpx"))
	gui.Run()
}
