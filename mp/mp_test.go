package mp

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestMacroDefine1(t *testing.T) {
	mp := New()
	file, err := os.Open("macro_define_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skips first MACRO line

	goal := []string{" ADD #(1,1) ", " #(1,2) #(1,1) ", " MACRO ", " TEST #(2,1) ", " ADD #(2,1) ", " MEND "}
	mp.macroDefine(scanner)
	macro := mp.macroDefinitiontable["M1"]
	for i := range goal {
		//fmt.Println(macro.instructions[i], goal[i])
		//fmt.Println(goal[i], macro.instructions[i])
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

func TestMacroDefine2(t *testing.T) {
	file, err := os.Open("macro_define_test2.asm")

	if err != nil {
		panic(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skips comment
	scanner.Scan() // skips first MACRO line

	mp := New()
	mp.macroDefine(scanner)
	macro := mp.macroDefinitiontable["P"]
	goal := []string{}
	goal = append(goal, "#(1,1) #(1,2) #(1,3) #(1,4) ")
	goal = append(goal, " MACRO ")
	goal = append(goal, " Q #(2,1) #(2,2) #(2,3) #(2,4) ")
	goal = append(goal, "#(2,1) #(2,2) #(1,3) #(1,4) ")
	goal = append(goal, " MACRO ")
	goal = append(goal, " R #(3,1) #(3,2) #(3,3) #(3,4) ")
	goal = append(goal, "#(3,1) #(2,2) #(3,2) #(1,4) ")
	goal = append(goal, "#(3,3) #(2,4) #(3,4) H ")
	goal = append(goal, " MEND ")
	goal = append(goal, "#(2,3) #(2,4) G H ")
	goal = append(goal, "E F G H ")
	goal = append(goal, " MEND ")

	for i := range goal {
		fmt.Println(macro.instructions[i], goal[i])
		if macro.instructions[i] != goal[i] {
			t.Fatalf("macro define is not working properly")
		}
	}
}
