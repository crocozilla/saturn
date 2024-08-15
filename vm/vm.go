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
	//OperationMode
	operation     shared.Operation
	memoryAddress uint16
	operations    map[shared.Operation]func(shared.Operands, shared.AddressMode)
	isRunning     bool
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.setupOperations()
	vm.stackInit()
	return vm
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

func (vm *VirtualMachine) stackPop() (uint16, error) {
	if vm.stackPointer == 0 {
		return 0, errors.New("empty stack")
	}

	address := vm.stackPointer + stackBase
	vm.stackPointer--

	return uint16(vm.memory[address]), nil
}

func (vm *VirtualMachine) Execute(instr shared.Instruction) {
	vm.programCounter++
	vm.operations[instr.Operation](instr.Operands, instr.AddressMode)
}

func (vm *VirtualMachine) ExecuteAll(program shared.Program) {
	vm.isRunning = true

	for vm.isRunning {
		currentInstruction := program[vm.programCounter]
		vm.Execute(currentInstruction)
	}

	vm.programCounter = 0
}
func extractAddressMode(instr shared.Instruction) shared.AddressMode {
	addressModeBits := int(instr.Operation) >> 4

	addressModes := map[uint16]shared.AddressMode{
		0b01_00: shared.DIRECT,
		0b10_00: shared.INDIRECT,
		0b11_00: shared.IMMEDIATE,
		0b01_10: shared.DIRECT_INDIRECT,
		0b10_01: shared.INDIRECT_DIRECT,
		0b01_11: shared.DIRECT_IMMEDIATE,
		0b10_11: shared.INDIRECT_IMMEDIATE,
	}

	mode, ok := addressModes[uint16(addressModeBits)]
	if !ok {
		panic("invalid address mode in instruction")
	}

	return mode
}

func extractOpCode(instr shared.Instruction) shared.Operation {
	return instr.Operation % 16
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
