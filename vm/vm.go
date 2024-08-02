package vm

import (
	"errors"
	"saturn/shared"
)

const (
	stackBase shared.Word = 2 //2 é definido no pdf do trabalho
)

type VirtualMachine struct {
	memory       [128]shared.Word
	pc           uint16
	stackPointer shared.Word
	accumulator  shared.Word
	//OperationMode
	operation    shared.Operation
	memoryAdress uint16
	operations   map[shared.Operation]func(shared.Operands)
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	vm.stackInit()
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

func (vm *VirtualMachine) stackInit() {
	// primeiro elemento da pilha é seu limite (definido no pdf)
	var stackLimit shared.Word = 10 //max elements
	vm.memory[stackBase] = stackLimit
	vm.stackPointer++
}

func (vm *VirtualMachine) stackPush(value shared.Word) error {
	stackLimit := vm.memory[stackBase]

	if vm.stackPointer <= stackLimit {
		pointer := vm.stackPointer + stackBase

		vm.memory[pointer] = value
		vm.stackPointer++
		return nil

	}

	vm.stackPointer = 0
	return errors.New("stack overflow")

}

func (vm *VirtualMachine) stackPop() (shared.Word, error) {
	vm.stackPointer--

	if vm.stackPointer > 0 { //cant pop first element (stackBase)
		pointer := vm.stackPointer + stackBase
		return vm.memory[pointer], nil
	}

	return 0, errors.New("empty stack")

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
// levar em consideracao modos de enderecamento
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
