package assembler

var pseudoInstructions map[string]bool = map[string]bool{
	"START": 	true,
	"END":		true,
	"INTDEF": 	true,
	"INTUSE": 	true,
	"CONST":  	true,
	"SPACE":  	true,
	"STACK":  	true,
}

