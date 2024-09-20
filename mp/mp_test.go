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
	mp := New()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skips comment
	scanner.Scan() // skips first MACRO line

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
	goal = append(goal, " MEND ")
	goal = append(goal, "E F G H ")
	goal = append(goal, " MEND ")

	for i := range goal {
		fmt.Println(macro.instructions[i], goal[i])
		if macro.instructions[i] != goal[i] {
			t.Fatalf("macro define is not working properly")
		}
	}
}

func TestMacroExpand(t *testing.T) {
	mp := New()
	file, err := os.Open("macro_expand_test.asm")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()          // skips first MACRO line
	mp.macroDefine(scanner) // defines C
	mp.macroDefine(scanner) // defines A
	mp.macroDefine(scanner) // defines B

	write_file, err := os.Create("macro_expansion")
	if err != nil {
		panic(err)
	}
	defer write_file.Close()
	defer os.Remove("macro_expansion")

	mp.macroExpand(" B TEST TEST2 ", write_file)
	write_file.Seek(0, 0)
	write_scanner := bufio.NewScanner(write_file)
	goal := []string{}
	goal = append(goal, " LOAD TEST ")
	goal = append(goal, " ADD TEST2 ")
	goal = append(goal, " ADD TEST2 ")
	goal = append(goal, " SUB TEST2 ")
	goal = append(goal, " SUB TEST2 ")

	for i := 0; write_scanner.Scan(); i++ {
		text := write_scanner.Text()
		fmt.Println(text, goal[i])
		if goal[i] != text {
			t.Fatalf("erro em macro expand")
		}

	}

}

func TestMacroPass(t *testing.T) {
	mp := New()
	file, err := os.Open("macro_pass_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mp.MacroPass(file)
}
