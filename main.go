package main

import (
	"saturn/assembler"
)

func main() {
	assembler.Run(
		"linker_test.asm",
		"linker_test_part2.asm")
}
