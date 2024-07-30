package assembler

type Instruction uint8

const (
	ADD    Instruction = 0b0010 // 2
	BR     Instruction = 0b0000 // 0
	BRNEG  Instruction = 0b0101 // 5
	BRPOS  Instruction = 0b0001 // 1
	BRZERO Instruction = 0b0100 // 4
	COPY   Instruction = 0b1101 // 13
	DIVIDE Instruction = 0b1010 // 10
	LOAD   Instruction = 0b0011 // 3
	MULT   Instruction = 0b1110 // 14
	READ   Instruction = 0b1100 // 12
	STOP   Instruction = 0b1011 // 11
	STORE  Instruction = 0b0111 // 7
	SUB    Instruction = 0b0110 // 6
	WRITE  Instruction = 0b1000 // 8
)

func (i Instruction) Check() bool {
	allowedInstructions := map[Instruction]bool{
		ADD:    true,
		BR:     true,
		BRNEG:  true,
		BRPOS:  true,
		BRZERO: true,
		COPY:   true,
		DIVIDE: true,
		LOAD:   true,
		MULT:   true,
		READ:   true,
		STOP:   true,
		STORE:  true,
		SUB:    true,
		WRITE:  true,
	}

	_, ok := allowedInstructions[i]

	return ok
}

func (i Instruction) String() {

}
