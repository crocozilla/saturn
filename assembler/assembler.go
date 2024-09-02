package assembler

import (
	"bufio"
	"errors"
	"os"
	"saturn/shared"
	"unicode"
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

	assembler.firstPass(file)
	assembler.secondPass(file)

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

func (assembler *Assembler) firstPass(file *os.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		// if operation is a pseudo-instruction, op2 is always empty
		label, operationString, op1, op2 := parseLine(line)

		if len(label) > 8 {
			panic("um símbolo excede o limite de caracteres (8).")
		}

		for i, v := range label {
			if i == 0 {
				if !unicode.IsLetter(v) {
					panic("primeiro caracter de um símbolo deve ser alfabético")
				}
			} else if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
				panic("símbolo deve apenas conter caracteres alfanuméricos")
			}
		}

		_, isPseudoInstruction := pseudoOpSizes[operationString]
		if isPseudoInstruction {
			instruction := operationString
			switch instruction {
			case "START":
			case "END":
				return
			case "INTDEF":
			case "INTUSE":
			case "CONST":
			case "SPACE":
			case "STACK":
			}

		} else {
			operation, err := getOpcode(operationString)
			if err != nil {
				panic("operação inválida.")
			}

			if len(label) != 0 {
				assembler.insertIntoSymbolTable(label)
			}

			assembler.locationCounter += shared.OpSizes[operation]
		}

	}

	panic("sem instrução \"end\"")

}

func (assembler *Assembler) secondPass(file *os.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		// if operation is a pseudo-instruction, op2 is always empty
		//_, operationString, op1, op2 := parseLine(line)

	}
}

func (assembler *Assembler) insertIntoSymbolTable(label string) {
	assembler.symbolTable[label] = assembler.locationCounter
}

// converts operand from string to its value
func getOperandValue(operand string) {
	panic("getOperandValue not implemented")
}
