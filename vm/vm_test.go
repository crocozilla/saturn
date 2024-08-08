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
	operands.First = 35
	instr.Operands = operands
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 50 {
		t.Fatalf(`add direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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

	//acc >= 0
	vm.accumulator = 0
	vm.programCounter = 0

	vm.Execute(instr)

	if shared.Word(vm.programCounter) != 0 {
		t.Fatalf(`brneg not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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

	//acc <= 0
	vm.accumulator = 0
	vm.programCounter = 0

	vm.Execute(instr)

	if vm.programCounter != 0 {
		t.Fatalf(`brpos not working, expected 0 on pc, got %v`, vm.programCounter)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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
		t.Fatalf(`br direct not working, expected 50 on pc, got %v`, vm.programCounter)
	}

	//acc != 0
	vm.accumulator = 1
	vm.programCounter = 0

	vm.Execute(instr)

	if vm.programCounter != 0 {
		t.Fatalf(`brzero not working, expected 0 on pc, got %v`, vm.programCounter)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
}

func TestCall(t *testing.T) {
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
		t.Fatalf(`call not working, expected 10 on stack, got %v`, vm.programCounter)
	}

	//direct

	if shared.Word(vm.programCounter) != vm.memory[operands.First] {
		t.Fatalf(`call direct not working, expected %v on pc, got %v`, vm.memory[operands.First], vm.programCounter)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
}

func TestCopy(t *testing.T) {
	t.Fatalf(`copy testing not implemented`)
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

	if vm.accumulator != 10/2 {
		t.Fatalf(`divide immediate not working, expected 5 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	operands.First = 35
	instr.Operands = operands
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 10/5 {
		t.Fatalf(`divide direct not working, expected 2 on acc, got %v`, vm.accumulator)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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
	operands.First = 35
	instr.Operands = operands
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 50 {
		t.Fatalf(`load direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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
	operands.First = 35
	instr.Operands = operands
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 10*5 {
		t.Fatalf(`mult direct not working, expected 50 on acc, got %v`, vm.accumulator)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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

	if vm.memory[operands.First] != vm.accumulator {
		t.Fatalf(`mult direct not working, expected 10 on memory, got %v`, vm.memory[operands.First])
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
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

	if vm.accumulator != 50-20 {
		t.Fatalf(`add immediate not working, expected 30 on acc, got %v`, vm.accumulator)
	}

	//direct
	vm.accumulator = 50
	vm.memory[35] = 50
	operands.First = 35
	instr.Operands = operands
	instr.AddressMode = shared.DIRECT

	vm.Execute(instr)

	if vm.accumulator != 0 {
		t.Fatalf(`add direct not working, expected 0 on acc, got %v`, vm.accumulator)
	}

	//indirect
	t.Fatalf(`indirect testing not implemented, other tests passed`)
}

func TestWrite(t *testing.T) {
	t.Fatalf(`write testing not implemented`)
}
