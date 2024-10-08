package assembler

import (
	"fmt"
	"math"
	"os"
	"saturn/shared"
	"strconv"
	"testing"
)

func TestFirstPass(t *testing.T) {
	file, err := os.Open("assembler_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	assembler := New()

	assembler.firstPass(file)

	//symbol table:
	label := "SIG"
	info, ok := assembler.symbolTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in symbol table")

	}
	if info.Address != 0 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "PROG"
	if _, ok := assembler.symbolTable[label]; ok {
		t.Fatalf(label + " should be in definition table")
	}

	//definition table:
	label = "PROG"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if info.Address != 0 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "UP"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")

	}
	if info.Address != 6 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "DOWN"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")

	}
	if info.Address != 7 {
		t.Fatalf("incorrect address on label " + label)
	}

	//use table:
	label = "LOOP"
	slice, ok := assembler.useTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")

	}
	if len(slice) == 0 {
		t.Fatalf("slice is empty on label " + label)
	}
	if len(slice) > 1 || slice[0] != 3 {
		t.Fatalf("incorrect slice in label " + label)
	}

	label = "X"
	slice, ok = assembler.useTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if len(slice) != 0 {
		t.Fatalf("incorrect slice in label " + label)
	}

	file2, err := os.Open("assembler_test_2.asm")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	assembler = New()

	assembler.firstPass(file2)

	// symbol table
	label = "X"
	info, ok = assembler.symbolTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in symbol table")

	}
	if info.Address != 14 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "LOC"
	info, ok = assembler.symbolTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in symbol table")
	}
	if info.Address != 12 {
		t.Fatalf("incorrect address on label " + label)
	}

	// definition table
	label = "PROG2"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if info.Address != 0 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "CAS"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if info.Address != 13 {
		t.Fatalf("incorrect address on label " + label)
	}

	label = "BRA"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if info.Address != 7 {
		fmt.Println(info.Address)
		t.Fatalf("incorrect address on label " + label)
	}

	label = "LOOP"
	info, ok = assembler.definitionTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in definition table")
	}
	if info.Address != 0 {
		t.Fatalf("incorrect address on label " + label)
	}

	//use table:
	label = "UP"
	slice, ok = assembler.useTable[label]
	if !ok {
		t.Fatalf("missing " + label + " in use table")
	}
	if len(slice) == 0 {
		t.Fatalf("incorrect slice in label " + label)
	}
	if slice[0] != 3 {
		t.Fatalf("incorrect slice in label " + label)
	}
}

func TestSecondPass(t *testing.T) {
	file, err := os.Open("assembler_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	assembler := New()

	assembler.firstPass(file)
	assembler.secondPass(file)

	opcode := shared.ADD
	operand1 := "1"
	operand2 := EMPTY

	assembler.addAddressModeToOpcode(&opcode, operand1, operand2)

	if shared.ExtractAddressMode(shared.Word(opcode)) != shared.DIRECT {
		t.Fatalf("incorrect address mode on second pass")
	}

	opcode = shared.ADD
	operand1 = "1,I"
	operand2 = EMPTY

	assembler.addAddressModeToOpcode(&opcode, operand1, operand2)

	if shared.ExtractAddressMode(shared.Word(opcode)) != shared.INDIRECT {
		t.Fatalf("incorrect address mode on second pass")
	}

	opcode = shared.ADD
	operand1 = "#1"
	operand2 = EMPTY

	assembler.addAddressModeToOpcode(&opcode, operand1, operand2)

	if shared.ExtractAddressMode(shared.Word(opcode)) != shared.IMMEDIATE {
		t.Fatalf("incorrect address mode on second pass")
	}

	opcode = shared.COPY
	operand1 = "1,I"
	operand2 = "#2"

	assembler.addAddressModeToOpcode(&opcode, operand1, operand2)

	if shared.ExtractAddressMode(shared.Word(opcode)) != shared.INDIRECT_IMMEDIATE {
		t.Fatalf("incorrect address mode on second pass")
	}

	os.Remove("FIRST.lst")
	os.Remove("FIRST.obj")
}

// missing literal tests
func TestGetOperandValue(t *testing.T) {
	iTest, err := getOperandValue("#64")
	if err != nil {
		panic(err)
	}
	if iTest != 64 {
		t.Fatalf("número imediato inválido")
	}
	hexTest, err := getOperandValue("H'F'")
	if err != nil {
		panic(err)
	}
	if hexTest != 15 {
		t.Fatalf("número hexadecimal inválido, esperava-se %v mas obteve-se %v", 15, hexTest)
	}
	hexTest, err = getOperandValue("H'32AF'")
	if err != nil {
		panic(err)
	}
	if hexTest != 12975 {
		t.Fatalf("número hexadecimal inválido, esperava-se %v mas obteve-se %v", 12975, hexTest)
	}
	//test, err = getOperandValue("@30")
	maxNumber := int(math.Pow(2.0, float64(shared.WordSize-1)) - 1.0)
	minNumber := int(-math.Pow(2.0, float64(shared.WordSize-1)))

	_, err = getOperandValue(strconv.Itoa(maxNumber + 1))
	if err == nil {
		t.Fatalf("número grande demais não gerou erro")
	}
	_, err = getOperandValue(strconv.Itoa(maxNumber))
	if err != nil {
		t.Fatalf("número máximo gerou erro incorretamente")
	}
	_, err = getOperandValue(strconv.Itoa(minNumber - 1))
	if err == nil {
		t.Fatalf("número pequeno demais não gerou erro")
	}
	_, err = getOperandValue(strconv.Itoa(minNumber))
	if err != nil {
		t.Fatalf("número mínimo negativo gerou erro incorretamente")
	}
	_, err = getOperandValue("-")
	if err == nil {
		t.Fatalf("número negativo vazio não gerou erro")
	}
}

func TestGetAddressMode(t *testing.T) {
	addressMode, err := getAddressMode("64")
	if err != nil {
		panic(err)
	}
	if addressMode != shared.DIRECT {
		t.Fatalf("addressMode deveria ser direto")
	}
	addressMode, err = getAddressMode("64,I")
	if err != nil {
		panic(err)
	}
	if addressMode != shared.INDIRECT {
		t.Fatalf("addressMode deveria ser indireto")
	}
	addressMode, err = getAddressMode("#64")
	if err != nil {
		panic(err)
	}
	if addressMode != shared.IMMEDIATE {
		t.Fatalf("addressMode deveria ser imediato")
	}
	_, err = getAddressMode("#64I")
	if err == nil {
		t.Fatalf("addressMode múltiplo foi permitido")
	}
	_, err = getAddressMode("64I")
	if err == nil {
		t.Fatalf("modo indireto incorreto foi permitido")
	}
}

func TestRemoveAddressMode(t *testing.T) {
	operand, _ := removeAddressMode("64,I")
	if operand != "64" {
		t.Fatalf("esperava-se operando 64, recebeu-se %v", operand)
	}
	operand, _ = removeAddressMode("#65")
	if operand != "65" {
		t.Fatalf("esperava-se operando 65, recebeu-se %v", operand)
	}
	operand, _ = removeAddressMode("#@H'3F'")
	if operand != "@H'3F'" {
		t.Fatalf("esperava-se operando @H'3F, recebeu-se %v", operand)
	}
}
