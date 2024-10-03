package main

import (
	"saturn/assembler"
	// "saturn/gui"
)

//"saturn/gui"

func main() {
	// program := ReadProgram("./program.txt")
	// gui.InsertProgram(program)
	// gui.Run()

	assembler.Run(
		"assembler/linker_test.asm",
		"assembler/linker_test_part2.asm")
}
