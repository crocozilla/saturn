package assembler

import (
	"bufio"
	"errors"
	"fmt"
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

type ObjCode []string

type Assembler struct {
	symbolTable     map[string]symbolInfo
	definitionTable map[string]symbolInfo
	useTable        map[string][]uint16
	locationCounter uint16
	lineCounter     uint16
	programName     string
	errors          []string
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
		assembler.lineCounter++
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		// if operation is a pseudo-instruction, op2 is always EMPTY
		label, operationString, op1, op2 := parseLine(line)
		if len(label) > 0 && validateSymbol(label) != nil {
			panic(validateSymbol(label))
		}
		op1SymbolErr := validateSymbol(op1)

		if _, ok := assembler.useTable[op1]; ok {
			assembler.useTable[op1] = append(assembler.useTable[op1], assembler.locationCounter+1)
		}
		if _, ok := assembler.useTable[op2]; ok {
			assembler.useTable[op2] = append(assembler.useTable[op2], assembler.locationCounter+2)
		}

		pseudoOpSize, isPseudoInstruction := pseudoOpSizes[operationString]
		if isPseudoInstruction {
			instruction := operationString
			switch instruction {
			case "START":
				if op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução start."))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
				if op1SymbolErr != nil {
					assembler.addError(errors.New("nome do programa inválido na pseudo instrução start."))
				}
				assembler.programName = op1
			case "END":
				if op1 != EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução end."))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
				return
			case "INTDEF":
				if op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução intdef."))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
				if op1SymbolErr == nil {
					// if a symbol is defined using intdef, it should be relocated from the symbolTable
					delete(assembler.symbolTable, op1)
					assembler.definitionTable[op1] = symbolInfo{assembler.locationCounter, ABSOLUTE}
				}
			case "INTUSE":
				if label == EMPTY || op1 != EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução intuse."))
				}
				assembler.useTable[label] = []uint16{}
			case "CONST":
				if label == EMPTY || op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução const."))
				}
				assembler.insertIntoProperTable(label)
			case "SPACE":
				if label == EMPTY || op1 != EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução space."))
				}
				assembler.insertIntoProperTable(label)
			case "STACK":
				if op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução stack."))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
			}
			assembler.locationCounter += pseudoOpSize
		} else {
			opcode, err := getOpcode(operationString)
			if err != nil {
				assembler.addError(errors.New("operação " + operationString + " é inválida."))
			}

			opSize := shared.OpSizes[opcode]
			sizeOneError := opSize == 1 && (op1 != EMPTY || op2 != EMPTY)
			sizeTwoError := opSize == 2 && (op1 == EMPTY || op2 != EMPTY)
			sizeThreeError := opSize == 3 && (op1 == EMPTY || op2 == EMPTY)
			invalidSyntax := sizeOneError || sizeTwoError || sizeThreeError
			if invalidSyntax {
				assembler.addError(errors.New("sintaxe inválida na operação " + operationString + "."))
			}

			if len(label) != 0 {
				assembler.insertIntoProperTable(label)
			}

			assembler.locationCounter += opSize
		}

	}

	assembler.addError(errors.New("sem instrução \"end\"."))

}

func (assembler *Assembler) secondPass(file *os.File) {
	// File rewind to origin and reset locationCount
	file.Seek(0, 0)
	assembler.locationCounter = 0
	// fmt.Println(assembler.symbolTable)

	if assembler.programName == EMPTY {
		assembler.addError(errors.New("programa sem nome"))
	}
	objFile, err := os.Create(assembler.programName + ".obj")
	if err != nil {
		panic(err)
	}
	lstFile, err := os.Create(assembler.programName + ".lst")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		_, operation, operand1, operand2 := parseLine(line)

		_, isPseudoInstruction := pseudoOpSizes[operation]
		if isPseudoInstruction {
			fmt.Println("Pseudo")
		} else {
			opCode, err := getOpcode(operation)
			if err != nil {
				panic(err)
			}

			var op1Value shared.Word
			if operand1 != EMPTY {
				if err := validateSymbol(operand1); err != nil {
					// Is a label
					info, okSymbol := assembler.symbolTable[operand1]

					addrUse, okUse := assembler.useTable[operand1]
					addrUse = addrUse

					addrDef, okDef := assembler.definitionTable[operand1]
					addrDef = addrDef

					if (okSymbol && okUse) || (okSymbol && okDef) || (okDef && okUse) {
						assembler.addError(errors.New("label " + operand1 + " defined in multiple tables"))
					}

					if okSymbol {

						op1Value = shared.Word(info.address)
					} else if okUse {
						// op1Value = shared.Word(addrUse)

					} else if okDef {
						// op1Value = shared.Word(addrDef.address)

					}

				} else {
					// Is a number
					op1Value, err = getOperandValue(operand1)
					if err != nil {
						panic(err)
					}
				}

			}

			var op2Value shared.Word
			if operand2 != EMPTY {
				if err := validateSymbol(operand2); err != nil {
					// Is a label
					info, ok := assembler.symbolTable[operand2]
					if !ok {
						assembler.addError(errors.New("not found label"))
					}

					op2Value = shared.Word(info.address)
				} else {
					// Is a number
					op2Value, err = getOperandValue(operand2)
					if err != nil {
						panic(err)
					}
				}
			}

			// Write a new line to obj file
			outputLine := fmt.Sprintf("%d %d %d \n", opCode, op1Value, op2Value)
			_, err = objFile.WriteString(outputLine)
			if err != nil {
				panic(err)
			}
		}
	}

	assembler.writeErrorsToLst(lstFile)
}

func (assembler *Assembler) addError(err error) {
	lineNumber := strconv.Itoa(int(assembler.lineCounter))
	errString := "erro na linha " + lineNumber + ": " + err.Error() + "\n"
	assembler.errors = append(assembler.errors, errString)
}

func (assembler *Assembler) writeErrorsToLst(lstFile *os.File) {
	if len(assembler.errors) == 0 {
		_, err := lstFile.WriteString("Nenhum erro detectado.\n")
		if err != nil {
			panic(err)
		}
		return
	}
	for _, errorString := range assembler.errors {
		_, err := lstFile.WriteString(errorString)
		if err != nil {
			panic(err)
		}
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

// checks for validity
func (assembler *Assembler) insertIntoProperTable(symbol string) {
	err := validateSymbol(symbol)
	if err != nil {
		panic(err)
	}

	// if its defined and it is its first use, set its address to current address
	info, ok := assembler.definitionTable[symbol]
	if ok && validateSymbol(symbol) == nil {
		if info.mode == ABSOLUTE {
			assembler.definitionTable[symbol] = symbolInfo{assembler.locationCounter, RELATIVE}
		}
	}

	_, ok = assembler.symbolTable[symbol]
	if ok {
		panic("símbolo " + symbol + " com múltiplas definições.")
	}
	assembler.symbolTable[symbol] = symbolInfo{assembler.locationCounter, RELATIVE}
}
