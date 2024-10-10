package vm

import (
	"errors"
	"saturn/shared"
)

const stackBase uint16 = 2 // 2 é definido no pdf do trabalho

type VirtualMachine struct {
	memory         [128]shared.Word
	programCounter uint16
	stackPointer   uint16
	accumulator    shared.Word
	operation      shared.Operation
	memoryAddress  uint16
	opImpls        map[shared.Operation]func(shared.Operands, shared.AddressMode)
	isRunning      bool
	stackLimit     uint16
	programBase    uint16
	programEnd     uint16
	io             struct {
		input  shared.Word
		output shared.Word
	}
}

func New(stackLimitArg uint16) *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	vm.stackInit()
	vm.isRunning = true
	vm.stackLimit = stackLimitArg
	vm.programBase = stackBase + vm.stackLimit + 1
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

func (vm *VirtualMachine) Output() uint16 {
	return uint16(vm.io.output)
}

func (vm *VirtualMachine) SetInput(data uint16) {
	vm.io.input = shared.Word(data)
}

func (vm *VirtualMachine) LoadProgram(program []shared.Word) {
	var i uint16

	for i = 0; i < uint16(len(program)); i++ {
		putAddress := vm.programBase + i

		if putAddress >= uint16(len(vm.memory)) {
			panic("the program exceeds the memory space")
		}

		vm.memory[putAddress] = program[i]
	}

	vm.programEnd = i + vm.programBase
}

func (vm *VirtualMachine) setupOperations() {
	vm.opImpls = map[shared.Operation]func(shared.Operands, shared.AddressMode){
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

func (vm *VirtualMachine) stackInit() {
	vm.memory[stackBase] = shared.Word(vm.stackLimit) // primeiro elemento da pilha é seu limite (definido no pdf)
}

func (vm *VirtualMachine) stackPush(value shared.Word) error {
	vm.stackPointer++

	if vm.stackPointer > vm.stackLimit {
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

	for i := 0; i < int(vm.programBase); i++ {
		vm.memory[i] = 0
	}

	for i := vm.programEnd; i < uint16(len(vm.memory)); i++ {
		vm.memory[i] = 0
	}

	vm.stackInit()
	vm.isRunning = true
}

func (vm *VirtualMachine) decodeInst() shared.Instruction {
	address := vm.programBase + vm.programCounter
	operationInfo := vm.memory[address]

	var operands shared.Operands

	operands.First = vm.memory[address+1]
	operands.Second = vm.memory[address+2] // might be trash, but when it is, it won`t be used by the instruction

	return shared.Instruction{
		AddressMode: shared.ExtractAddressMode(operationInfo),
		Operation:   shared.ExtractOpCode(operationInfo),
		Operands:    operands,
	}
}

func (vm *VirtualMachine) Execute() {
	instr := vm.decodeInst()

	vm.operation = instr.Operation
	vm.programCounter += shared.OpSizes[instr.Operation]

	vm.opImpls[instr.Operation](instr.Operands, instr.AddressMode)
}

func (vm *VirtualMachine) ExecuteAll() {
	vm.Reset()

	for vm.isRunning {
		vm.Execute()
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
	switch mode {
	case shared.DIRECT:
		vm.memory[operands.First] = vm.io.input
	case shared.DIRECT_INDIRECT:
		vm.memory[vm.memoryAddress] = vm.io.input
	default:
		panic("incorrect address mode on READ operation")
	}
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
	switch mode {
	case shared.IMMEDIATE:
		vm.io.output = operands.First
	case shared.DIRECT:
		vm.io.output = vm.memory[operands.First]
	case shared.INDIRECT:
		vm.io.output = vm.memory[vm.memoryAddress]
	default:
		panic("incorrect address mode on WRITE operation")
	}
}

func (vm *VirtualMachine) inj(operands shared.Operands, mode shared.AddressMode) {
	if mode == shared.IMMEDIATE {
		vm.memoryAddress = uint16(operands.First)
	} else {
		panic("incorrect address mode on INJECT operation")
	}
}
