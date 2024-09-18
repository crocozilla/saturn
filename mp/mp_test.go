package mp

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestMacroDefine(t *testing.T) {
	mp := New()
	file, err := os.Open("macro_define_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skips first MACRO line

	goal := []string{" ADD #1 ", " #2 #1 ", " MACRO ", " TEST ", " ADD @2 ", " MEND "}
	mp.macroDefine(scanner)
	macro := mp.macroDefinitiontable["M1"]
	for i := range goal {
		fmt.Println(macro.instructions[i], goal[i])
		if macro.instructions[i] != goal[i] {
			t.Fatalf("macro define is not working properly")
		}
	}
	if mp.macroDefinitiontable["M1"].numberOfParameters != 2 {
		t.Fatalf("incorrect number of parameters in macro M1")
	}

	if mp.macroDefinitiontable["TEST"].numberOfParameters != 0 {
		t.Fatalf("incorrect number of parameters in macro TEST")
	}
}
