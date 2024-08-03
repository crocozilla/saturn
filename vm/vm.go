package vm

import (
	"errors"
	"saturn/shared"
)

const stackBase uint16 = 2 // 2 é definido no pdf do trabalho

type addressMode uint8

const (
	DIRECT addressMode = iota
	INDIRECT_01
	INDIRECT_10
	INDIRECT_11
	IMMEDIATE
)

type VirtualMachine struct {
	memory         [128]shared.Word
	programCounter uint16
	stackPointer   uint16
	accumulator    shared.Word
	//OperationMode
	operation    shared.Operation
	memoryAdress uint16
	operations   map[shared.Operation]func(shared.Operands, addressMode)
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	vm.stackInit()
	return vm
}

func (vm *VirtualMachine) setupOperations() {
	// Provavelmente não está funcionando pois a shared.Operation tem modos de endereçamento acoplados no OPCODE
	vm.operations = map[shared.Operation]func(shared.Operands, addressMode){
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
	var stackLimit uint16 = 10                     // max elements
	vm.memory[stackBase] = shared.Word(stackLimit) // primeiro elemento da pilha é seu limite (definido no pdf)
}

func (vm *VirtualMachine) stackPush(value shared.Word) error {
	vm.stackPointer++
	stackLimit := uint16(vm.memory[stackBase])

	if vm.stackPointer > stackLimit {
		vm.stackPointer = 0
		return errors.New("stack overflow")
	}

	address := vm.stackPointer + stackBase
	vm.memory[address] = value

	return nil
}

func (vm *VirtualMachine) stackPop() (shared.Word, error) {
	if vm.stackPointer == 0 {
		return 0, errors.New("empty stack")
	}

	address := vm.stackPointer + stackBase
	vm.stackPointer--

	return vm.memory[address], nil
}

func (vm *VirtualMachine) Execute(instr shared.Instruction) {
	addressMode := extractAddressMode(instr)
	vm.operations[instr.Operation](instr.Operands, addressMode)
}

func (vm *VirtualMachine) ExecuteAll(program shared.Program) {
	instr := program[vm.programCounter]

	for instr.Operation != shared.STOP {
		vm.Execute(instr)
		vm.programCounter++
		instr = program[vm.programCounter]
	}
}

func extractAddressMode(instr shared.Instruction) addressMode {
	bitMask := 0b0000000001110000

	decidingBits := (bitMask & int(instr.Operation)) >> 4

	addressModes := map[uint16]addressMode{
		0b000: DIRECT,
		0b001: INDIRECT_01,
		0b010: INDIRECT_10,
		0b011: INDIRECT_11,
		0b100: IMMEDIATE,
	}

	mode, ok := addressModes[uint16(decidingBits)]
	if !ok {
		panic("invalid address mode in instruction")
	}

	return mode
}

// -- Operations

func (vm *VirtualMachine) add(operands shared.Operands, mode addressMode) {
	vm.accumulator = vm.accumulator + operands.First
}

func (vm *VirtualMachine) br(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) brneg(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) brpos(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) brzero(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) call(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) copy(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) divide(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) load(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) mult(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) read(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) ret(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) stop(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) store(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) sub(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) write(operands shared.Operands, mode addressMode) {
	panic("not implemented")
}
