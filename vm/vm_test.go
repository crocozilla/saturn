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

	vm.Execute(instr)

	if vm.accumulator != 20 {
		t.Fatalf(`add immediate not working, expected 20 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 0
	vm.memory[35] = 50
	instr.Operands.First = 35
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 50 {
		t.Fatalf(`add direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	vm.accumulator = 0
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.accumulator != 10 {
		t.Fatalf(`add indirect not working, expected 10 on acc, got %v`, vm.accumulator)
	}
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

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 50 {
		t.Fatalf(`br direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 10 {
		t.Fatalf(`br indirect not working, expected 10 on pc, got %v`, vm.programCounter)
	}
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

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 50 {
		t.Fatalf(`brneg direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 10 {
		t.Fatalf(`brneg indirect not working, expected 10 on pc, got %v`, vm.programCounter)
	}

	//acc >= 0
	vm.accumulator = 0
	vm.programCounter = 0

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 0 {
		t.Fatalf(`brneg not working, expected 0 on pc, got %v`, vm.programCounter)
	}
}

func TestBrpos(t *testing.T) {
	vm := New()

	//direct
	vm.accumulator = 1
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BRPOS,
		Operands:    operands}

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 50 {
		t.Fatalf(`brpos direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 10 {
		t.Fatalf(`brpos indirect not working, expected 10 on pc, got %v`, vm.programCounter)
	}

	//acc <= 0
	vm.accumulator = 0
	vm.programCounter = 0

	vm.Execute(instr)

	if vm.programCounter != 0 {
		t.Fatalf(`brpos not working, expected 0 on pc, got %v`, vm.programCounter)
	}

}
func TestBrzero(t *testing.T) {
	vm := New()

	//direct
	operands := shared.Operands{First: 35, Second: 0}
	vm.memory[35] = 50
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.BRZERO,
		Operands:    operands}

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 50 {
		t.Fatalf(`brzero direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 10 {
		t.Fatalf(`brzero indirect not working, expected 10 on pc, got %v`, vm.programCounter)
	}

	//acc != 0
	vm.accumulator = 1
	vm.programCounter = 0

	vm.Execute(instr)

	if vm.programCounter != 0 {
		t.Fatalf(`brzero not working, expected 0 on pc, got %v`, vm.programCounter)
	}
}

func TestCall(t *testing.T) {
	t.Fatalf(`call testing not defined yet because program counter logic may change`)

	vm := New()

	//stack tests
	vm.programCounter = 10
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

	if uint16(result) != 10 {
		t.Fatalf(`call not working, expected 10 on stack, got %v`, result)
	}

	//direct
	if vm.programCounter != 50 {
		t.Fatalf(`call direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	vm.memoryAddress = 40
	vm.memory[40] = 20
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.programCounter != 20 {
		t.Fatalf(`call indirect not working, expected 20 on pc, got %v`, vm.programCounter)
	}
}

func TestCopy(t *testing.T) {
	vm := New()

	//both direct
	operands := shared.Operands{First: 30, Second: 40}
	instr := shared.Instruction{
		AddressMode: shared.DIRECT,
		Operation:   shared.COPY,
		Operands:    operands}
	vm.memory[40] = 20

	vm.Execute(instr)

	if vm.memory[30] != 20 {
		t.Fatalf(`copy direct not working, expected 20 on memory, got %v`, vm.memory[30])
	}

	//direct_indirect
	vm.memoryAddress = 50
	vm.memory[50] = 60
	instr.Operands.First = 70
	instr.AddressMode = shared.DIRECT_INDIRECT

	vm.Execute(instr)

	if vm.memory[70] != 60 {
		t.Fatalf(`copy direct_indirect not working, expected 60 on memory, got %v`, vm.memory[70])
	}

	//direct_immediate
	instr.Operands.First = 30
	instr.Operands.Second = 40
	instr.AddressMode = shared.DIRECT_IMMEDIATE

	vm.Execute(instr)

	if vm.memory[30] != 40 {
		t.Fatalf(`copy direct_immediate not working, expected 40 on memory, got %v`, vm.memory[30])
	}

	//indirect_direct
	vm.memoryAddress = 90
	vm.memory[40] = 20
	instr.Operands.Second = 40
	instr.AddressMode = shared.INDIRECT_DIRECT

	vm.Execute(instr)

	if vm.memory[90] != 20 {
		t.Fatalf(`copy indirect_direct not working, expected 20 on memory, got %v`, vm.memory[90])
	}

	//both indirect doesnt make sense to test

	//indirect_immediate
	vm.memoryAddress = 50
	instr.Operands.Second = 100
	instr.AddressMode = shared.INDIRECT_IMMEDIATE

	vm.Execute(instr)

	if vm.memory[50] != 100 {
		t.Fatalf(`copy indirect_immediate not working, expected 100 on memory, got %v`, vm.memory[50])
	}
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

	vm.Execute(instr)

	if vm.accumulator != 5 {
		t.Fatalf(`divide immediate not working, expected 5 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	instr.Operands.First = 35
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 2 {
		t.Fatalf(`divide direct not working, expected 2 on acc, got %v`, vm.accumulator)
	}

	//indirect
	vm.accumulator = 20
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.accumulator != 2 {
		t.Fatalf(`divide indirect not working, expected 2 on acc, got %v`, vm.accumulator)
	}
}

func TestLoad(t *testing.T) {
	vm := New()

	//immediate
	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.LOAD,
		Operands:    operands}

	vm.Execute(instr)

	if vm.accumulator != 20 {
		t.Fatalf(`load immediate not working, expected 20 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 0
	vm.memory[35] = 50
	instr.Operands.First = 35
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 50 {
		t.Fatalf(`load direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.accumulator != 10 {
		t.Fatalf(`load indirect not working, expected 10 on acc, got %v`, vm.accumulator)
	}
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

	vm.Execute(instr)

	if vm.accumulator != 20 {
		t.Fatalf(`mult immediate not working, expected 20 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	instr.Operands.First = 35
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 50 {
		t.Fatalf(`mult direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	vm.accumulator = 2
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.accumulator != 20 {
		t.Fatalf(`mult indirect not working, expected 20 on acc, got %v`, vm.accumulator)
	}
}

func TestRead(t *testing.T) {
	t.Fatalf(`read testing not implemented`)
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

	vm.Execute(instr)

	if vm.programCounter != uint16(testAdress) {
		t.Fatalf(`ret not working, expected %v on pc, got %v`, testAdress, vm.programCounter)
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

	vm.Execute(instr)

	if vm.memory[operands.First] != 10 {
		t.Fatalf(`store direct not working, expected 10 on memory, got %v`, vm.memory[operands.First])
	}

	//indirect
	vm.accumulator = 25
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.memory[vm.memoryAddress] != 25 {
		t.Fatalf(`store indirect not working, expected 25 on memory, got %v`, vm.memory[vm.memoryAddress])
	}
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

	vm.Execute(instr)

	if vm.accumulator != 30 {
		t.Fatalf(`sub immediate not working, expected 30 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 50
	vm.memory[35] = 40
	instr.Operands.First = 35
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 10 {
		t.Fatalf(`sub direct not working, expected 10 on acc, got %v`, vm.accumulator)
	}

	//indirect
	vm.accumulator = 20
	vm.memory[30] = 10
	vm.memoryAddress = 30
	instr.AddressMode = shared.INDIRECT

	vm.Execute(instr)

	if vm.accumulator != 10 {
		t.Fatalf(`sub indirect not working, expected 10 on acc, got %v`, vm.accumulator)
	}
}

func TestWrite(t *testing.T) {
	t.Fatalf(`write testing not implemented`)
}

func TestInj(t *testing.T) {
	vm := New()

	operands := shared.Operands{First: 20, Second: 0}
	instr := shared.Instruction{
		AddressMode: shared.IMMEDIATE,
		Operation:   shared.INJ,
		Operands:    operands}

	vm.Execute(instr)

	if vm.memoryAddress != 20 {
		t.Fatalf(`inj not working, expected 20 on memory register, got %v`, vm.memoryAddress)
	}
}
