package main

import (
	"saturn/assembler"
	"saturn/linker"
)

func main() {
	//assembler.Run(
	//	"linker/linker_test.asm",
	//	"linker/linker_test_part2.asm")
	linker.Run(assembler.Run("linker/linker_test.asm", "linker/linker_test_part2.asm"))
}
