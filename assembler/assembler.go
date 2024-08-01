package assembler

import (
	"errors"
	"saturn/shared"
)

var sourceCodePath string

func Run(filePath string) (program shared.Program) {
	sourceCodePath = filePath

	scanLines(sourceCodePath, firstStep)
	scanLines(sourceCodePath, secondStep)

	// to-do: assemble program
	return shared.Program{}
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

func firstStep(line string) {
	// to-do
}

func secondStep(line string) {
	// to-do
}
