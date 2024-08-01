package vm

import "saturn/shared"

type VirtualMachine struct {
	memory [128]shared.Word
	pc     uint16
	//stack pointer
	accumulator shared.Word
	//OperationMode
	operation    shared.Operation // talvez deveria ser tipo shared.operation
	memoryAdress uint16
	operations   map[shared.Operation]func(shared.Operands)
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	return vm
}

func (vm *VirtualMachine) setupOperations() {
	vm.operations = map[shared.Operation]func(shared.Operands){
		shared.ADD:    vm.add,
		shared.BR:     vm.br,
		shared.BRNEG:  vm.brneg,
		shared.BRPOS:  vm.brpos,
		shared.BRZERO: vm.brzero,
		shared.CALL:   vm.call,
		shared.COPY:   vm.copy,
		shared.DIVIDE: vm.divide,
		shared.LOAD:   vm.load,
		shared.MULT:   vm.mult,
		shared.READ:   vm.read,
		shared.RET:    vm.ret,
		shared.STOP:   vm.stop,
		shared.STORE:  vm.store,
		shared.SUB:    vm.sub,
		shared.WRITE:  vm.write,
	}
}

func (vm *VirtualMachine) Execute(instr shared.Instruction) {
	vm.operations[instr.Operation](instr.Operands)
}

func (vm *VirtualMachine) ExecuteAll(program shared.Program) {
	for _, instr := range program {
		vm.Execute(instr)
	}
}

// -- Operations

func (vm *VirtualMachine) add(operands shared.Operands) {
	vm.accumulator = vm.accumulator + operands.First
}

func (vm *VirtualMachine) br(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) brneg(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) brpos(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) brzero(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) call(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) copy(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) divide(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) load(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) mult(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) read(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) ret(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) stop(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) store(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) sub(operands shared.Operands) {
	panic("not implemented")
}

func (vm *VirtualMachine) write(operands shared.Operands) {
	panic("not implemented")
}
