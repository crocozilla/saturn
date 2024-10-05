package main

import (
	"saturn/assembler"
)

func main() {
	//assembler.Run(
	//	"linker/linker_test.asm",
	//	"linker/linker_test_part2.asm")
	assembler.Run("linker/linker_test.asm", "linker/linker_test_part2.asm")
}
