package shared

import "fmt"

type Program []Instruction

type Instruction struct {
	Operation Operation
	Operands  Operands
}

type Operands struct {
	First  byte
	Second byte
}

type Operation byte

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
)

// string -> Operation -> func

// Campo de operação -> (Instr (ADD, SUB...) ou PseudoInstr (CONST, SPACE...))
// func (i Instruction) Check() bool {
// 	// allowedInstructions := map[Instruction]bool{
// 	// 	ADD:    true,
// 	// 	BR:     true,
// 	// 	BRNEG:  true,
// 	// 	BRPOS:  true,
// 	// 	BRZERO: true,
// 	// 	COPY:   true,
// 	// 	DIVIDE: true,
// 	// 	LOAD:   true,
// 	// 	MULT:   true,
// 	// 	READ:   true,
// 	// 	STOP:   true,
// 	// 	STORE:  true,
// 	// 	SUB:    true,
// 	// 	WRITE:  true,
// 	// }

// 	return ok
	// missing call and ret
// }

func (i Instruction) String() string {
	return fmt.Sprintf("<(%d) [%d, %d]>", i.Operation, i.Operands.First, i.Operands.Second)
}
