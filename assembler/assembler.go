package assembler

import (
	"bufio"
	"errors"
	"os"
	"saturn/shared"
	"strconv"
	"unicode"
)

var sourceCodePath string

type ObjCode []string

type Assembler struct {
	symbolTable     map[string]shared.Word
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
				value, err := getOperandValue(op1)
				if err != nil {
					panic(err)
				}
				assembler.symbolTable[label] = value
			case "SPACE":
			case "STACK":
			}

		} else {
			operation, err := getOpcode(operationString)
			if err != nil {
				panic("operação inválida.")
			}

			if len(label) != 0 {
				assembler.symbolTable[label] = shared.Word(assembler.locationCounter)
			}

			assembler.locationCounter += shared.OpSizes[operation]
		}

	}

	panic("sem instrução \"end\"")

}

func (assembler *Assembler) secondPass(file *os.File) {
	os.Create("result.obj")

	// File rewind to origin and reset locationCount
	file.Seek(0, 0)
	assembler.locationCounter = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}
	os.Create("result.obj")

	// File rewind to origin and reset locationCount
	file.Seek(0, 0)
	assembler.locationCounter = 0

		_, operation, operand1, operand2 := parseLine(line)

		opCode, err := getOpcode(operation)
		if err != nil {
			panic(err)
		}

		if _, ok := assembler.symbolTable[operand1]; !ok {
			assembler.symbolTable[operand1] = shared.Word(assembler.locationCounter)
		}

		if _, ok := assembler.symbolTable[operand2]; !ok {
			assembler.symbolTable[operand2] = shared.Word(assembler.locationCounter)
		}

		assembler.locationCounter++


	}
}

func getOperandValue(operand string) (shared.Word, error) {
	if len(operand) == 0 {
		return shared.Word(0), errors.New("operando vazio usado em getOperandValue")
	}

	var value int64
	var err error
	operand, err = RemoveAdressFromOperand(operand)
	if err != nil {
		return shared.Word(0), err
	}
	apostrophe := byte('\'')
	isHexadecimal := operand[0] == 'H' && len(operand) > 3
	isLiteral := operand[0] == '@' && len(operand) > 1
	switch {
	case isHexadecimal:
		if operand[1] != apostrophe && operand[len(operand)-1] != apostrophe {
			return 0, errors.New("faltando apostrofos em número hexadecimal")
		}
		hexString := operand[1 : len(operand)-1]
		value, err = strconv.ParseInt(hexString, 16, shared.WordSize)
		if err != nil {
			return 0, errors.New("número hexadecimal inválido")
		}

	case isLiteral:
		value, err = strconv.ParseInt(operand[1:], 10, shared.WordSize)
		if err != nil {
			return 0, errors.New("literal decimal inválido")
		}
	default:
		value, err = strconv.ParseInt(operand, 10, shared.WordSize)
		if err != nil {
			return 0, errors.New("número não reconhecido")
		}
	}

	return shared.Word(value), nil
}

func RemoveAdressFromOperand(operand string) (string, error) {
	if operand[0] == '#' && len(operand) > 1 {
		operand = operand[1:]
	} else if len(operand) > 2 && operand[len(operand)-2] == ',' && operand[len(operand)-1] == 'I' {
		operand = operand[0 : len(operand)-2]
	} else if _, err := strconv.Atoi(operand); err != nil {
		return "", errors.New("operando inválido")
	}

	return operand, nil
}
