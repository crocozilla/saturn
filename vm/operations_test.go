package vm

import (
	"saturn/shared"
	"testing"
)

func TestAdd(t *testing.T) {
	vm := New()

	//immediate
	operands := shared.Operands{20, 0}
	instr := shared.Instruction{shared.IMMEDIATE, shared.ADD, operands}
	var expected shared.Word = 20
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 0
	vm.memory[35] = 50
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = 50
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	vm.accumulator = 0

	//indirect
	panic("indirect testing not implemented")
}

func TestBr(t *testing.T) {
	panic("not implemented")
}

func TestBrneg(t *testing.T) {
	panic("not implemented")
}

func TestBrzero(t *testing.T) {
	panic("not implemented")
}

func TestCall(t *testing.T) {
	panic("not implemented")
}

func TestCopy(t *testing.T) {
	panic("not implemented")
}

func TestDivide(t *testing.T) {
	vm := New()

	//immediate
	vm.accumulator = 10
	operands := shared.Operands{2, 0}
	instr := shared.Instruction{shared.IMMEDIATE, shared.ADD, operands}
	var expected shared.Word = 10 / 2
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add immediate not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//direct
	vm.accumulator = 10
	vm.memory[35] = 5
	operands.First = 35
	instr.AddressMode = shared.DIRECT
	expected = 10 / 5
	vm.Execute(instr)

	if vm.accumulator != expected {
		t.Fatalf(`add direct not working, expected %v on acc, got %v`, expected, vm.accumulator)
	}

	//indirect
	panic("indirect testing not implemented")
}

func TestLoad(t *testing.T) {
	panic("not implemented")
}

func TestMult(t *testing.T) {
	panic("not implemented")
}

func TestRead(t *testing.T) {
	panic("not implemented")
}

func TestRet(t *testing.T) {
	panic("not implemented")
}

func TestStore(t *testing.T) {
	panic("not implemented")
}

func TestSub(t *testing.T) {
	panic("not implemented")
}

func TestWrite(t *testing.T) {
	panic("not implemented")
}
