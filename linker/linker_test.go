package linker

import (
	"saturn/assembler"
	"testing"
)

func TestRun(t *testing.T) {
	_ = Run(assembler.Run("linker_test.asm", "linker_test_part2.asm"))
}
