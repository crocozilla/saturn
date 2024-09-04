package assembler

// numbers arent set properly yet
var pseudoOpSizes map[string]uint16 = map[string]uint16{
	"START":  0,
	"END":    0,
	"INTDEF": 1,
	"INTUSE": 1,
	"CONST":  0,
	"SPACE":  1,
	"STACK":  1,
}

/*
func treatPseudoInstruction(instruction string, operand string) {
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
}*/
