package linker

import (
	"saturn/assembler"
	"testing"
)

func TestRun(t *testing.T) {
	_ = Run(assembler.Run("linker_test.asm", "linker_test_part2.asm"))
	_ = Run(assembler.Run("linker_test_3.asm", "linker_test_3.asm"))
	// todo: compare first run with MAIN_test goal
	// and second run with TESTE3_test goal
}
