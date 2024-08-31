package assembler

import (
	"errors"
	"saturn/shared"
	"bufio"
	"os"
)

var sourceCodePath string

type Assembler struct{
	symbolTable map[string]shared.Word // dont know if should be shared.Word
	locationCounter int
	end bool
}

func New() *Assembler{
	assembler := new(Assembler)
	return assembler
}

// writes to file program.txt as its output
func Run(filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	assembler := New()

	assembler.firstStep(file)
	assembler.secondStep(file)

	// to-do: assemble program
}

func Check(token string) (shared.Operation, error) {
	allowedInstructions := map[string]shared.Operation{
		"ADD":    2,
		"BR":     0,
		"BRNEG":  5,
		"BRPOS":  1,
		"BRZERO": 4,
		"CALL":   15,
		"COPY":   13,
		"DIVIDE": 10,
		"LOAD":   3,
		"MULT":   14,
		"READ":   12,
		"RET":    16,
		"STOP":   11,
		"STORE":  7,
		"SUB":    6,
		"WRITE":  8,
	}

	if opCode, ok := allowedInstructions[token]; ok {
		return opCode, nil
	} else {
		return 99, errors.New("token not found")
	}
}

func (assembler *Assembler) firstStep(file *os.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		line := scanner.Text()
		if (len(line) > 80){
			panic("linha muito longa. Não deve haver mais de 80 caracteres numa linha.")
		}

		// whole line is a comment
		if(line[0] == '*'){
			continue
		}

		//label, operation, op1, op2 := parseLine(line)
/*
		// pseudoInstructions is a map defined in pseudo_instruction.go
		if(_, ok := pseudoInstructions[operation]; ok){
			treatPseudoInstruction(operation)
			if(assembler.end == true){
				return
			}
		}

		if len(label) == 0{
			// add to symbol table
			 
		}
*/
	}

}

func (assembler *Assembler) secondStep(file *os.File) {
	// to-do
	panic("secondStep not implemented")
}

func treatPseudoInstruction(assembler Assembler, pseudoInstruction string){
	switch (pseudoInstruction){
	case "START":
	case "END":
		assembler.end = true
	case "INTDEF":
	case "INTUSE":
	case "CONST":
	case "SPACE":
	case "STACK":
	}
}
