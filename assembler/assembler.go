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
	lstLineCounter  uint16
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
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução start"))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
				if op1SymbolErr != nil {
					assembler.addError(errors.New("nome do programa inválido na pseudo instrução start"))
				}
				assembler.programName = op1
			case "END":
				if op1 != EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução end"))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
				return
			case "INTDEF":
				if op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução intdef"))
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
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução intuse"))
				}
				assembler.useTable[label] = []uint16{}
			case "CONST":
				if label == EMPTY || op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução const"))
				}
				assembler.insertIntoProperTable(label)
			case "SPACE":
				if label == EMPTY || op1 != EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução space"))
				}
				assembler.insertIntoProperTable(label)
			case "STACK":
				if op1 == EMPTY || op2 != EMPTY {
					assembler.addError(errors.New("sintaxe inválida na pseudo instrução stack"))
				}
				if label != EMPTY {
					assembler.insertIntoProperTable(label)
				}
			}
			assembler.locationCounter += pseudoOpSize
		} else {
			opcode, err := getOpcode(operationString)
			if err != nil {
				assembler.addError(errors.New("operação " + operationString + " é inválida"))
			}

			opSize := shared.OpSizes[opcode]
			sizeOneError := opSize == 1 && (op1 != EMPTY || op2 != EMPTY)
			sizeTwoError := opSize == 2 && (op1 == EMPTY || op2 != EMPTY)
			sizeThreeError := opSize == 3 && (op1 == EMPTY || op2 == EMPTY)
			invalidSyntax := sizeOneError || sizeTwoError || sizeThreeError
			if invalidSyntax {
				assembler.addError(errors.New("sintaxe inválida na operação " + operationString))
			}

			if len(label) != 0 {
				assembler.insertIntoProperTable(label)
			}

			assembler.locationCounter += opSize
		}

	}

	assembler.addError(errors.New("sem instrução \"end\""))

}

func (assembler *Assembler) secondPass(file *os.File) {
	fmt.Println(assembler.useTable)
	// File rewind to origin and reset locationCount
	file.Seek(0, 0)
	assembler.locationCounter = 0

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
	assembler.lineCounter = 0
	assembler.lstLineCounter = 1
	var op1Value, op2Value shared.Word
	var op1Mode, op2Mode byte
	var zeroValuedByte byte
	var opCode shared.Operation

	for scanner.Scan() {
		assembleLine := false
		assembler.lineCounter++

		line, isComment := readLine(scanner)
		if isComment {
			continue
		}

		// method 'assembleLine' needs these values zeroed if unused
		op1Mode = zeroValuedByte
		op2Mode = zeroValuedByte

		_, operation, operand1, operand2 := parseLine(line)

		fmt.Printf("%s %s %s\n", operation, operand1, operand2)

		opSize, isPseudoInstruction := pseudoOpSizes[operation]
		if isPseudoInstruction {
			switch operation {
			case "CONST":
				if operand1 != EMPTY {
					op1Value, op1Mode = assembler.getOperandValueAndMode(operand1)
				}
				assembleLine = true
			default:
				assembleLine = false
			}

		} else {
			assembleLine = true
			opCode, err = getOpcode(operation)
			if err != nil {
				panic(err)
			}
			// redefines opSize
			opSize = shared.OpSizes[opCode]

			if operand1 != EMPTY {
				op1Value, op1Mode = assembler.getOperandValueAndMode(operand1)
			}

			if operand2 != EMPTY {
				op2Value, op2Mode = assembler.getOperandValueAndMode(operand2)
			}

		}

		assembler.addAddressModeToOpcode(&opCode, operand1, operand2)

		if assembleLine {
			assembler.assembleLine(objFile, lstFile, isPseudoInstruction, opCode, op1Value, op1Mode, op2Value, op2Mode)
		}

		assembler.locationCounter += opSize

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

func (assembler *Assembler) addAddressModeToOpcode(opCode *shared.Operation, operand1 string, operand2 string) {
	if operand1 != EMPTY {
		op1AddressMode, err := getAddressMode(operand1)
		if err != nil {
			assembler.addError(err)
		}
		if op1AddressMode == shared.DIRECT {
			*opCode += 0b01_00 << 5
		} else if op1AddressMode == shared.INDIRECT {
			*opCode += 0b10_00 << 5
		} else if op1AddressMode == shared.IMMEDIATE {
			*opCode += 0b11_00 << 5
		}
	}

	if operand2 != EMPTY {
		op2AddressMode, err := getAddressMode(operand2)
		if err != nil {
			assembler.addError(err)
		}
		if op2AddressMode == shared.DIRECT {
			*opCode += 0b00_01 << 5
		} else if op2AddressMode == shared.INDIRECT {
			*opCode += 0b00_10 << 5
		} else if op2AddressMode == shared.IMMEDIATE {
			*opCode += 0b00_11 << 5
		}
	}

}

// checks if mode is unset to see if operands are being used, ignores them if needed
func (assembler *Assembler) assembleLine(objFile *os.File, lstFile *os.File, isPseudoInstruction bool,
	opCode shared.Operation, op1Value shared.Word, op1Mode byte, op2Value shared.Word, op2Mode byte) {
	var objLine string
	var lstLine string = fmt.Sprintf("%02d ", assembler.locationCounter)
	var zeroValuedByte byte
	smallPadding := "    "
	padding := "     "

	if !isPseudoInstruction {
		opCodeString := fmt.Sprintf("%02d ", opCode)
		objLine += opCodeString
		lstLine += opCodeString
	} else {
		objLine += smallPadding
		lstLine += smallPadding
	}

	if op1Mode != zeroValuedByte {
		op1String := fmt.Sprintf("%02d %c ", op1Value, op1Mode)
		objLine += op1String
		lstLine += op1String
	} else {
		lstLine += padding
	}
	if op2Mode != zeroValuedByte {
		op2String := fmt.Sprintf("%02d %c ", op2Value, op2Mode)
		objLine += op2String
		lstLine += op2String
	} else {
		lstLine += padding
	}

	lstLine += fmt.Sprintf("%02d %02d", assembler.lstLineCounter, assembler.lineCounter)

	objLine += "\n"
	lstLine += "\n"

	_, err := objFile.WriteString(objLine)
	if err != nil {
		panic(err)
	}

	_, err = lstFile.WriteString(lstLine)
	if err != nil {
		panic(err)
	}

	assembler.lstLineCounter++
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
			return 0, errors.New("número " + operand + " não reconhecido")
		}
	}

	return shared.Word(value), nil
}

// assumes operand is not empty
func (assembler *Assembler) getOperandValueAndMode(operand string) (value shared.Word, mode byte) {
	if err := validateSymbol(operand); err == nil {
		// Is a label
		infoSym, okSym := assembler.symbolTable[operand]
		infoDef, okDef := assembler.definitionTable[operand]

		_, okUse := assembler.useTable[operand]

		// TODO: Check if can be in multiple tables
		// if (okSymbol && okUse) || (okSymbol && okDef) || (okDef && okUse) {
		// 	assembler.addError(errors.New("label " + operand1 + " defined in multiple tables"))
		// }

		if okSym {
			value = shared.Word(infoSym.address)
			mode = infoSym.mode
		} else if okUse {
			// ?
			value = 0
			mode = 'A'
		} else if okDef {
			value = shared.Word(infoDef.address)
			mode = infoDef.mode
		}
	} else {
		// Is a number
		value, err = getOperandValue(operand)
		if err != nil {
			panic(err)
		}
		mode = 'A'
	}

	return value, mode
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
		assembler.addError(err)
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
		assembler.addError(errors.New("símbolo " + symbol + " com múltiplas definições."))
	}
	assembler.symbolTable[symbol] = symbolInfo{assembler.locationCounter, RELATIVE}
}
