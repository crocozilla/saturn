package shared

import "fmt"

type Word int16

type Program []Instruction

type BinProgram []BinInstruction

type AddressMode uint8

const (
	DIRECT AddressMode = iota
	INDIRECT
	IMMEDIATE
	DIRECT_INDIRECT
	DIRECT_IMMEDIATE
	INDIRECT_DIRECT
	INDIRECT_IMMEDIATE
)

type Instruction struct {
	AddressMode AddressMode
	Operation   Operation
	Operands    Operands
}

type Operands struct {
	First  Word
	Second Word
}

type Operation Word

const (
	ADD    Operation = 2
	BR     Operation = 0
	BRNEG  Operation = 5
	BRPOS  Operation = 1
	BRZERO Operation = 4
	CALL   Operation = 15
	COPY   Operation = 13
	DIVIDE Operation = 10
	LOAD   Operation = 3
	MULT   Operation = 14
	READ   Operation = 12
	RET    Operation = 16
	STOP   Operation = 11
	STORE  Operation = 7
	SUB    Operation = 6
	WRITE  Operation = 8
	INJ    Operation = 9
)

// TESTAR TUDO
type BinInstruction [3]Word

// Binary Instruction to Instruction
func Btoi(b BinInstruction) Instruction {
	return Instruction{
		AddressMode: extractAddressMode(b[0]),
		Operation:   extractOpCode(b[0]),
		Operands:    Operands{First: b[1], Second: b[2]},
	}
}

// Binary Program to Program
func Btop(bp BinProgram) Program {
	p := make(Program, len(bp))

	for i, v := range bp {
		p[i] = Btoi(v)
	}

	return p
}

func extractAddressMode(operation Word) AddressMode {
	addressModeBits := int(operation) >> 4

	addressModes := map[uint16]AddressMode{
		0b01_00: DIRECT,
		0b10_00: INDIRECT,
		0b11_00: IMMEDIATE,
		0b01_10: DIRECT_INDIRECT,
		0b10_01: INDIRECT_DIRECT,
		0b01_11: DIRECT_IMMEDIATE,
		0b10_11: INDIRECT_IMMEDIATE,
	}

	mode, ok := addressModes[uint16(addressModeBits)]
	if !ok {
		panic("invalid address mode in instruction")
	}

	return mode
}

func extractOpCode(operation Word) Operation {
	return Operation(operation % 16)
}

func (i Instruction) String() string {
	return fmt.Sprintf("<(%d) [%d, %d]>", i.Operation, i.Operands.First, i.Operands.Second)
}
