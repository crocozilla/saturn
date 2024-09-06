package assembler

import (
	"bufio"
	"errors"
	"os"
	"saturn/shared"
	"strconv"
	"unicode"
)

const (
	RELATIVE = 'R'
	ABSOLUTE = 'A'
)

type symbolInfo struct {
	address uint16
	mode    byte
}

type Assembler struct {
	symbolTable     map[string]symbolInfo
	definitionTable map[string]symbolInfo
	useTable        map[string][]uint16
	locationCounter uint16
	programName     string
}

func New() *Assembler {
	assembler := new(Assembler)
	assembler.symbolTable = map[string]symbolInfo{}
	assembler.definitionTable = map[string]symbolInfo{}
	assembler.useTable = map[string][]uint16{}
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

		// if operation is a pseudo-instruction, op2 is always EMPTY
		label, operationString, op1, op2 := parseLine(line)

		op1SymbolErr := validateSymbol(op1)
		op2SymbolErr := validateSymbol(op2)

		if _, ok := assembler.definitionTable[op1]; ok {
			assembler.useTable[op1] = append(assembler.useTable[op1], assembler.locationCounter+1)
		}
		if _, ok := assembler.definitionTable[op2]; ok {
			assembler.useTable[op2] = append(assembler.useTable[op2], assembler.locationCounter+2)
		}

		pseudoOpSize, isPseudoInstruction := pseudoOpSizes[operationString]
		if isPseudoInstruction {
			instruction := operationString
			switch instruction {
			case "START":
				if op1 == EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução start.")
				}
				if label != EMPTY {
					assembler.insertIntoSymbolTable(label, RELATIVE)
				}
				if op1SymbolErr != nil {
					panic("nome do programa inválido na pseudo instrução start.")
				}
				assembler.programName = op1
			case "END":
				if op1 != EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução end.")
				}
				if label != EMPTY {
					assembler.insertIntoSymbolTable(label, RELATIVE)
				}
				return
			case "INTDEF":
				if op1 == EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução intdef.")
				}
				if label != EMPTY {
					assembler.insertIntoSymbolTable(label, RELATIVE)
				}
				if op1SymbolErr != nil {
					// if a symbol is defined using intdef, it should be relocated from the symbolTable
					if _, ok := assembler.symbolTable[op1]; ok {
						delete(assembler.symbolTable, op1)
					}
					assembler.definitionTable[op1] = symbolInfo{assembler.locationCounter + 1, RELATIVE}
				}
			case "INTUSE":
				if label == EMPTY || op1 != EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução intuse.")
				}
				assembler.useTable[label] = []uint16{}
			case "CONST":
				if label == EMPTY || op1 == EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução const.")
				}

			case "SPACE":
				if label == EMPTY || op1 != EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução space.")
				}
				assembler.insertIntoSymbolTable(label, RELATIVE)
			case "STACK":
				if op1 == EMPTY || op2 != EMPTY {
					panic("sintaxe inválida na pseudo instrução stack.")
				}
				if label != EMPTY {
					assembler.insertIntoSymbolTable(label, RELATIVE)
				}
			}
			assembler.locationCounter += pseudoOpSize
		} else {
			opcode, err := getOpcode(operationString)
			if err != nil {
				panic("operação " + operationString + " é inválida.")
			}

			opSize := shared.OpSizes[opcode]
			sizeOneError := opSize == 1 && (op1 != EMPTY || op2 != EMPTY)
			sizeTwoError := opSize == 2 && (op1 == EMPTY || op2 != EMPTY)
			sizeThreeError := opSize == 3 && (op1 == EMPTY || op2 == EMPTY)
			invalidSyntax := sizeOneError || sizeTwoError || sizeThreeError
			if invalidSyntax {
				panic("sintaxe inválida na operação " + operationString + ".")
			}

			if len(label) != 0 {
				assembler.insertIntoSymbolTable(label, RELATIVE)
			}

			assembler.locationCounter += opSize
		}

	}

	panic("sem instrução \"end\".")

}

func (assembler *Assembler) secondPass(file *os.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		//_, operationString, op1, op2 := parseLine(line)

	}
}

func getOperandValue(operand string) (shared.Word, error) {
	if len(operand) == 0 {
		return shared.Word(0), errors.New("operando vazio usado em getOperandValue")
	}

	var value int64
	var err error
	operand, err = removeAddressMode(operand)
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
		hexString := operand[2 : len(operand)-1]
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

func removeAddressMode(operand string) (string, error) {
	addressMode, err := getAddressMode(operand)
	if err != nil {
		return EMPTY, err
	}
	if addressMode == shared.IMMEDIATE {
		operand = operand[1:]
	} else if addressMode == shared.INDIRECT {
		operand = operand[0 : len(operand)-2]
	}

	return operand, nil
}

func getAddressMode(operand string) (shared.AddressMode, error) {
	if operand == EMPTY {
		return shared.DIRECT, errors.New("operando vazio em getAddressMode")
	}
	if operand[0] == '#' && len(operand) > 1 && operand[len(operand)-1] == 'I' {
		return shared.DIRECT, errors.New("operando com múltiplos endereçamentos em getAddressMode")
	}
	// 									 note the "!="
	if len(operand) > 2 && operand[len(operand)-2] != ',' && operand[len(operand)-1] == 'I' {
		return shared.DIRECT, errors.New("operando indireto inválido")
	}
	if operand[0] == '#' && len(operand) > 1 {
		return shared.IMMEDIATE, nil
	} else if len(operand) > 2 && operand[len(operand)-2] == ',' && operand[len(operand)-1] == 'I' {
		return shared.INDIRECT, nil
	}

	return shared.DIRECT, nil
}

func validateSymbol(symbol string) error {
	if symbol == EMPTY {
		return errors.New("símbolo vazio em validateSymbol")
	}
	if len(symbol) > 8 {
		return errors.New("símbolo " + symbol + " excede o limite de caracteres (8).")
	}

	for i, v := range symbol {
		if i == 0 {
			if !unicode.IsLetter(v) {
				return errors.New("primeiro caracter de um símbolo deve ser alfabético")
			}
		} else if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
			return errors.New("símbolo deve apenas conter caracteres alfanuméricos")
		}
	}

	return nil
}

// assumes symbol is not empty, checks for validity
func (assembler *Assembler) insertIntoSymbolTable(symbol string, mode byte) {
	err := validateSymbol(symbol)
	if err != nil {
		panic(err)
	}

	_, ok := assembler.symbolTable[symbol]
	if ok {
		panic("símbolo " + symbol + " com múltiplas definições.")
	}
	assembler.symbolTable[symbol] = symbolInfo{assembler.locationCounter, mode}
}
