package vm

type Instr uint16

type InstrParameters struct {
	first  uint8
	second uint8
}

const (
	ADD    = 0b0010 // 2
	BR     = 0b0000 // 0
	BRNEG  = 0b0101 // 5
	BRPOS  = 0b0001 // 1
	BRZERO = 0b0100 // 4
	COPY   = 0b1101 // 13
	DIVIDE = 0b1010 // 10
	LOAD   = 0b0011 // 3
	MULT   = 0b1110 // 14
	READ   = 0b1100 // 12
	STOP   = 0b1011 // 11
	STORE  = 0b0111 // 7
	SUB    = 0b0110 // 6
	WRITE  = 0b1000 // 8
)

type VirtualMachine struct {
	accumulator uint32
	memory      [128]byte
	operations  map[Instr]func(InstrParameters)
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	return vm
}

func (vm *VirtualMachine) setupOperations() {
	vm.operations = map[Instr]func(InstrParameters){
		ADD: vm.add,
	}
}

func (vm *VirtualMachine) Execute(instruction Instr, parameters InstrParameters) {
	vm.operations[instruction](parameters)
}

func (vm *VirtualMachine) add(parameters InstrParameters) {
	
}
