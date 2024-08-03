package vm

import (
	"saturn/shared"
	"testing"
)

func Test_extractAddressMode(t *testing.T) {
	type info struct {
		expected addressMode
		code     uint16
	}

	cases := []info{
		{expected: IMMEDIATE, code: 0b0000000001000000},
		{expected: INDIRECT_01, code: 0b0000000000010000},
		{expected: INDIRECT_10, code: 0b0000000000100000},
		{expected: INDIRECT_11, code: 0b0000000000110000},
		{expected: DIRECT, code: 0b0000000000000000},
	}

	for _, info := range cases {
		instr := shared.Instruction{
			Operation: shared.Operation(info.code),
			Operands: shared.Operands{
				First:  0,
				Second: 0,
			},
		}

		mode := extractAddressMode(instr)

		if mode != info.expected {
			t.Fatalf("wrong result, expected: %d, got: %d", info.expected, mode)
		}
	}
}
