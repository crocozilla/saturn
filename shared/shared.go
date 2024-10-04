package shared

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Word int16

const WordSize = 16

type Program []Instruction

type BinProgram []BinInstruction

type AddressMode uint8

type SymbolInfo struct {
	Address uint16
	Mode    byte
}

const (
	DIRECT AddressMode = iota
	INDIRECT
	IMMEDIATE
	DIRECT_INDIRECT
	DIRECT_IMMEDIATE
	INDIRECT_DIRECT
	INDIRECT_IMMEDIATE
	UNUSED
)

const (
	RELATIVE = 'R'
	ABSOLUTE = 'A'
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

var OpSizes map[Operation]uint16 = map[Operation]uint16{
	ADD:    2,
	BR:     2,
	BRNEG:  2,
	BRPOS:  2,
	BRZERO: 2,
	CALL:   2,
	COPY:   3,
	DIVIDE: 2,
	LOAD:   2,
	MULT:   2,
	READ:   2,
	RET:    1,
	STOP:   1,
	STORE:  2,
	SUB:    2,
	WRITE:  2,
	INJ:    2,
}

// TESTAR TUDO
type BinInstruction [3]Word

// Binary Instruction to Instruction
func Btoi(b BinInstruction) Instruction {
	return Instruction{
		AddressMode: ExtractAddressMode(b[0]),
		Operation:   ExtractOpCode(b[0]),
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

func ExtractAddressMode(operation Word) AddressMode {
	addressModeBits := int(operation) >> 5

	addressModes := map[uint16]AddressMode{
		0b01_00: DIRECT,
		0b10_00: INDIRECT,
		0b11_00: IMMEDIATE,
		0b01_10: DIRECT_INDIRECT,
		0b10_01: INDIRECT_DIRECT,
		0b01_11: DIRECT_IMMEDIATE,
		0b10_11: INDIRECT_IMMEDIATE,
		0b00_00: UNUSED,
	}

	mode, ok := addressModes[uint16(addressModeBits)]
	if !ok {
		panic("invalid address mode in instruction")
	}

	return mode
}

func ExtractOpCode(operation Word) Operation {
	return Operation(operation & 0b0000000_00_00_11111)
}

func (i Instruction) String() string {
	return fmt.Sprintf("%d<(%d) [%d, %d]>", i.AddressMode, i.Operation, i.Operands.First, i.Operands.Second)
}

func CreateBuildFile(fileName string) (*os.File, error) {
	// if is test
	if strings.HasSuffix(os.Args[0], ".test") {
		buildPath := filepath.Join("..", "build")
		exists := directoryExists(buildPath)
		if !exists {
			os.MkdirAll(buildPath, 0777)
		}
		path := filepath.Join(buildPath, fileName)
		file, err := os.Create(path)
		return file, err
	} else {
		exists := directoryExists("build")
		if !exists {
			os.MkdirAll("build", 0777)
		}
		file, err := os.Create(filepath.Join("build", fileName))
		return file, err
	}
}

func OpenBuildFile(fileName string) (*os.File, error) {
	// if is test
	if strings.HasSuffix(os.Args[0], ".test") {
		buildPath := filepath.Join("..", "build")
		path := filepath.Join(buildPath, fileName)
		file, err := os.Open(path)
		return file, err

	} else {
		file, err := os.Open(filepath.Join("build", fileName))
		return file, err
	}
}

// exists returns whether the given directory exists
// panics on wrong kind of error
func directoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}
