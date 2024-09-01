package assembler

import (
	"bufio"
	"errors"
	"os"
	"saturn/shared"
)

var sourceCodePath string

type Assembler struct {
	symbolTable     map[string]uint16
	locationCounter uint16
}

func New() *Assembler {
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

func getOpcode(token string) (shared.Operation, error) {
	allowedInstructions := map[string]shared.Operation{
		"ADD":    shared.ADD,
		"BR":     shared.BR,
		"BRNEG":  shared.BRNEG,
		"BRPOS":  shared.BRPOS,
		"BRZERO": shared.BRZERO,
		"CALL":   shared.CALL,
		"COPY":   shared.COPY,
		"DIVIDE": shared.DIVIDE,
		"LOAD":   shared.LOAD,
		"MULT":   shared.MULT,
		"READ":   shared.READ,
		"RET":    shared.RET,
		"STOP":   shared.STOP,
		"STORE":  shared.STORE,
		"SUB":    shared.SUB,
		"WRITE":  shared.WRITE,
	}

	if opCode, ok := allowedInstructions[token]; ok {
		return opCode, nil
	} else {
		return 99, errors.New("token not found")
	}
}

func (assembler *Assembler) firstStep(file *os.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 80 {
			panic("linha muito longa. Não deve haver mais de 80 caracteres numa linha.")
		}

		// whole line is a comment
		if line[0] == '*' {
			continue
		}

		// if operation is a pseudo-instruction, op2 is always empty
		label, operationString, op1, op2 := parseLine(line)

		_, isPseudoInstruction := pseudoOpSizes[operationString]
		if isPseudoInstruction {
			treatPseudoInstruction(operationString, op1)

		} else {
			operation, err := getOpcode(operationString)
			if err != nil {
				panic("invalid operation")
			}

			if len(label) != 0 {
				assembler.insertIntoSymbolTable(label)
			}

			assembler.locationCounter += shared.OpSizes[operation]
		}

	}

	panic("no end instruction.")

}

func (assembler *Assembler) secondStep(file *os.File) {
	// to-do
	panic("secondStep not implemented")
}

func (assembler *Assembler) insertIntoSymbolTable(label string) {
	assembler.symbolTable[label] = assembler.locationCounter
}

// converts operand from string to its value
//func getOperandValue(operand string)  {}
