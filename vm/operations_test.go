package vm

import (
	"saturn/shared"
	"testing"
)

func TestAdd(t *testing.T) {
	vm := New()

	//immediate
	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.ADD,
		Operands:    operands}
	var expected shared.Word = vm.accumulator + operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 0
	vm.memory[35] = 50
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = vm.accumulator + vm.memory[operands.First]
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestBr(t *testing.T) {
	vm := New()

	//direct
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BR,
		Operands:    operands}
	var expected shared.Word = vm.memory[operands.First]
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`br direct not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestBrneg(t *testing.T) {
	vm := New()

	//direct
	vm.accumulator = -1
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BRNEG,
		Operands:    operands}
	var expected shared.Word = vm.memory[operands.First]
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`brneg direct not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//acc >= 0
	vm.accumulator = 0
	vm.programCounter = 0
	expected = shared.Word(vm.programCounter)
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`brneg not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}
	//indirect
	panic("indirect testing not implemented")
}

func TestBrpos(t *testing.T) {
	vm := New()

	//direct
	vm.accumulator = 1
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BRNEG,
		Operands:    operands}
	var expected shared.Word = vm.memory[operands.First]
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`brpos direct not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//acc <= 0
	vm.accumulator = 0
	vm.programCounter = 0
	expected = shared.Word(vm.programCounter)
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`brpos not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//indirect
	panic("indirect testing not implemented")
}
func TestBrzero(t *testing.T) {
	vm := New()

	//direct
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BRNEG,
		Operands:    operands}
	var expected shared.Word = vm.memory[operands.First]
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`br direct not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//acc != 0
	vm.accumulator = 1
	vm.programCounter = 0
	expected = shared.Word(vm.programCounter)
	vm.Execute(instr)

	if shared.Word(vm.programCounter) != expected {
		t.Fatalf(`brzero not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestCall(t *testing.T) {
	vm := New()

	//stack tests
	var expected uint16 = 10
	vm.programCounter = expected
	operands := shared.Operands{First: 35, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.CALL,
		Operands:    operands}
	vm.memory[35] = 50
	vm.Execute(instr)

	result, err := vm.stackPop()
	if err != nil {
		t.Fatalf(`call not working, stack is empty (should have program counter)`)
	}

	if uint16(result) != expected {
		t.Fatalf(`call not working, expected %v on stack, got %v`, expected, vm.programCounter)
	}

	//direct
	expected2 := vm.memory[operands.First]

	if vm.programCounter != uint16(expected2) {
		t.Fatalf(`call direct not working, expected %v on pc, got %v`, expected2, vm.programCounter)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestCopy(t *testing.T) {
	panic("not implemented")
}

func TestDivide(t *testing.T) {
	vm := New()

	//immediate
	vm.accumulator = 10
	operands := shared.Operands{First: 2, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.DIVIDE,
		Operands:    operands}
	var expected shared.Word = vm.accumulator / operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`divide immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = vm.accumulator / vm.memory[operands.First]
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`divide direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestLoad(t *testing.T) {
	vm := New()

	//immediate
	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.LOAD,
		Operands:    operands}
	var expected shared.Word = operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`load immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 0
	vm.memory[35] = 50
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = vm.memory[operands.First]
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`load direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestMult(t *testing.T) {
	vm := New()

	//immediate
	vm.accumulator = 10
	operands := shared.Operands{First: 2, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.MULT,
		Operands:    operands}
	var expected shared.Word = vm.accumulator * operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`mult immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = vm.accumulator * vm.memory[operands.First]
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`mult direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestRead(t *testing.T) {
	panic("not implemented")
}

func TestRet(t *testing.T) {
	vm := New()

	var testAdress shared.Word = 20
	vm.stackPush(testAdress)

	operands := shared.Operands{First: 0, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.RET,
		Operands:    operands}
	var expected uint16 = uint16(testAdress)
	vm.Execute(instr)

	if vm.programCounter != expected {
		t.Fatalf(`ret not working, expected %v on pc, got %v`, expected, vm.programCounter)
	}
}

func TestStore(t *testing.T) {
	vm := New()

	//direct
	vm.accumulator = 10
	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.STORE,
		Operands:    operands}
	var expected shared.Word = vm.accumulator
	vm.Execute(instr)

	if vm.memory[operands.First] != expected {
		t.Fatalf(`mult direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestSub(t *testing.T) {
	vm := New()

	//immediate
	vm.accumulator = 50
	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.SUB,
		Operands:    operands}
	var expected shared.Word = vm.accumulator - operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 50
	vm.memory[35] = 50
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = vm.accumulator - operands.First
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestWrite(t *testing.T) {
	panic("not implemented")
}
