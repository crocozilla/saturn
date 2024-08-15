package vm

import (
	"errors"
	"saturn/shared"
)

const stackBase uint16 = 2 // 2 é definido no pdf do trabalho
const stackLimit uint16 = 10
const programBase uint16 = stackBase + stackLimit + 1

type VirtualMachine struct {
	memory         [128]shared.Word
	programCounter uint16
	stackPointer   uint16
	accumulator    shared.Word
	operation      shared.Operation
	memoryAddress  uint16
	operations     map[shared.Operation]func(shared.Operands, shared.AddressMode)
	opSizes        map[shared.Operation]uint16
	isRunning      bool
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	vm.setupOpSizes()
	vm.stackInit()
	vm.isRunning = true
	return vm
}

func (vm *VirtualMachine) Memory() [128]shared.Word {
	return vm.memory
}

func (vm *VirtualMachine) PC() uint16 {
	return vm.programCounter
}

func (vm *VirtualMachine) SP() uint16 {
	return vm.stackPointer
}

func (vm *VirtualMachine) Accumulator() shared.Word {
	return vm.accumulator
}

func (vm *VirtualMachine) Operation() shared.Operation {
	return vm.operation
}

func (vm *VirtualMachine) MemoryAddress() uint16 {
	return vm.memoryAddress
}

func (vm *VirtualMachine) IsRunning() bool {
	return vm.isRunning
}

// provavelmente serão removidas, só para checkpoint
func (vm *VirtualMachine) TurnOn() {
	vm.isRunning = true
}

func (vm *VirtualMachine) TurnOff() {
	vm.isRunning = false
}

func (vm *VirtualMachine) InsertProgram(program []shared.Word) {

	for i := uint16(0); i < uint16(len(program)); i++ {
		vm.memory[programBase+i] = program[i]
	}
}

func (vm *VirtualMachine) setupOperations() {
	vm.operations = map[shared.Operation]func(shared.Operands, shared.AddressMode){
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
		shared.INJ:    vm.inj,
	}
}

func (vm *VirtualMachine) setupOpSizes() {
	vm.opSizes = map[shared.Operation]uint16{
		shared.ADD:    2,
		shared.BR:     2,
		shared.BRNEG:  2,
		shared.BRPOS:  2,
		shared.BRZERO: 2,
		shared.CALL:   2,
		shared.COPY:   3,
		shared.DIVIDE: 2,
		shared.LOAD:   2,
		shared.MULT:   2,
		shared.READ:   2,
		shared.RET:    1,
		shared.STOP:   1,
		shared.STORE:  2,
		shared.SUB:    2,
		shared.WRITE:  2,
		shared.INJ:    2,
	}
}

func (vm *VirtualMachine) stackInit() {
	//var stackLimit uint16 = 10                     // max elements
	vm.memory[stackBase] = shared.Word(stackLimit) // primeiro elemento da pilha é seu limite (definido no pdf)
}

func (vm *VirtualMachine) stackPush(value shared.Word) error {
	vm.stackPointer++
	//	stackLimit := uint16(vm.memory[stackBase])

	if vm.stackPointer > stackLimit {
		vm.stackPointer = 0
		return errors.New("stack overflow")
	}

	address := vm.stackPointer + stackBase
	vm.memory[address] = value

	return nil
}

func (vm *VirtualMachine) stackPop() (uint16, error) {
	if vm.stackPointer == 0 {
		return 0, errors.New("empty stack")
	}

	address := vm.stackPointer + stackBase
	vm.stackPointer--

	return uint16(vm.memory[address]), nil
}

func (vm *VirtualMachine) Reset() {
	vm.programCounter = 0
	vm.accumulator = 0
	vm.operation = 0
	vm.memoryAddress = 0
	vm.stackPointer = 0

	for i := range vm.memory {
		vm.memory[i] = 0
	}

	vm.stackInit()
	vm.isRunning = true
}

func (vm *VirtualMachine) structureInstruction(pc uint16) shared.Instruction {
	address := pc + programBase
	operation_info := vm.memory[address]

	instr := shared.Instruction{
		AddressMode: shared.ExtractAddressMode(operation_info),
		Operation:   shared.ExtractOpCode(operation_info),
		Operands:    shared.Operands{First: vm.memory[address+1], Second: vm.memory[address+2]}, // might be trash, but when it is, it won`t be used by the instruction
	}

	return instr
}

func (vm *VirtualMachine) Execute(pc uint16) {
	instr := vm.structureInstruction(pc)

	vm.operation = instr.Operation
	vm.programCounter += vm.opSizes[instr.Operation]

	vm.operations[instr.Operation](instr.Operands, instr.AddressMode)
}

func (vm *VirtualMachine) ExecuteAll(program shared.Program) {
	vm.Reset()

	for vm.isRunning {
		vm.Execute(vm.programCounter)
	}
}

// -- Operations

func (vm *VirtualMachine) add(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.IMMEDIATE:
		vm.accumulator += operands.First

	case shared.DIRECT:
		vm.accumulator += vm.memory[operands.First]

	case shared.INDIRECT:
		vm.accumulator += vm.memory[vm.memoryAddress]

	default:
		panic("incorrect address mode on ADD operation")
	}
}

func (vm *VirtualMachine) br(operands shared.Operands, mode shared.AddressMode) {
	var targetAddress uint16

	switch mode {
	case shared.DIRECT:
		targetAddress = uint16(vm.memory[operands.First])

	case shared.INDIRECT:
		targetAddress = uint16(vm.memory[vm.memoryAddress])

	default:
		panic("incorrect address mode on BR operation")
	}

	vm.programCounter = targetAddress
}

func (vm *VirtualMachine) brneg(operands shared.Operands, mode shared.AddressMode) {
	if vm.accumulator < 0 {
		vm.br(operands, mode)
	}
}

func (vm *VirtualMachine) brpos(operands shared.Operands, mode shared.AddressMode) {
	if vm.accumulator > 0 {
		vm.br(operands, mode)
	}
}

func (vm *VirtualMachine) brzero(operands shared.Operands, mode shared.AddressMode) {
	if vm.accumulator == 0 {
		vm.br(operands, mode)
	}
}

func (vm *VirtualMachine) call(operands shared.Operands, mode shared.AddressMode) {

	err := vm.stackPush(shared.Word(vm.programCounter))
	if err != nil {
		panic(err)
	}

	switch mode {
	case shared.DIRECT:
		vm.programCounter = uint16(vm.memory[operands.First])

	case shared.INDIRECT:
		vm.programCounter = uint16(vm.memory[vm.memoryAddress])

	default:
		panic("incorrect address mode on CALL operation")
	}
}

func (vm *VirtualMachine) copy(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.DIRECT:
		vm.memory[operands.First] = vm.memory[operands.Second]

	case shared.DIRECT_IMMEDIATE:
		vm.memory[operands.First] = operands.Second

	case shared.DIRECT_INDIRECT:
		vm.memory[operands.First] = vm.memory[vm.memoryAddress]

	case shared.INDIRECT:
		break

	case shared.INDIRECT_IMMEDIATE:
		vm.memory[vm.memoryAddress] = operands.Second

	case shared.INDIRECT_DIRECT:
		vm.memory[vm.memoryAddress] = vm.memory[operands.Second]

	default:
		panic("incorrect address mode on COPY operation")
	}
}

func (vm *VirtualMachine) divide(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.IMMEDIATE:
		vm.accumulator = vm.accumulator / operands.First

	case shared.DIRECT:
		vm.accumulator = vm.accumulator / vm.memory[operands.First]

	case shared.INDIRECT:
		vm.accumulator = vm.accumulator / vm.memory[vm.memoryAddress]

	default:
		panic("incorrect address mode on DIVIDE operation")
	}
}

func (vm *VirtualMachine) load(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.IMMEDIATE:
		vm.accumulator = operands.First

	case shared.DIRECT:
		vm.accumulator = vm.memory[operands.First]

	case shared.INDIRECT:
		vm.accumulator = vm.memory[vm.memoryAddress]

	default:
		panic("incorrect address mode on LOAD operation")
	}
}

func (vm *VirtualMachine) mult(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.IMMEDIATE:
		vm.accumulator = vm.accumulator * operands.First

	case shared.DIRECT:
		vm.accumulator = vm.accumulator * vm.memory[operands.First]

	case shared.INDIRECT:
		vm.accumulator = vm.accumulator * vm.memory[vm.memoryAddress]

	default:
		panic("incorrect address mode on MULT operation")
	}
}

func (vm *VirtualMachine) read(operands shared.Operands, mode shared.AddressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) ret(operands shared.Operands, mode shared.AddressMode) {
	var err error
	vm.programCounter, err = vm.stackPop()

	if err != nil {
		panic(err)
	}
}

func (vm *VirtualMachine) stop(operands shared.Operands, mode shared.AddressMode) {
	vm.isRunning = false
}

func (vm *VirtualMachine) store(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.DIRECT:
		vm.memory[operands.First] = vm.accumulator

	case shared.INDIRECT:
		vm.memory[vm.memoryAddress] = vm.accumulator

	default:
		panic("incorrect address mode on STORE operation")
	}
}

func (vm *VirtualMachine) sub(operands shared.Operands, mode shared.AddressMode) {
	switch mode {
	case shared.IMMEDIATE:
		vm.accumulator = vm.accumulator - operands.First

	case shared.DIRECT:
		vm.accumulator = vm.accumulator - vm.memory[operands.First]

	case shared.INDIRECT:
		vm.accumulator = vm.accumulator - vm.memory[vm.memoryAddress]

	default:
		panic("incorrect address mode on SUB operation")
	}
}

func (vm *VirtualMachine) write(operands shared.Operands, mode shared.AddressMode) {
	panic("not implemented")
}

func (vm *VirtualMachine) inj(operands shared.Operands, mode shared.AddressMode) {
	if mode == shared.IMMEDIATE {
		vm.memoryAddress = uint16(operands.First)
	} else {
		panic("incorrect address mode on INJECT operation")
	}
}
