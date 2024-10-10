package main

import (
	"bufio"
	"saturn/shared"
	"strconv"
	"strings"
)

func ReadProgram(fileName string) []shared.Word {
	file, err := shared.OpenBuildFile(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var program []shared.Word

	for scanner.Scan() {
		curr_instr := strings.Split(scanner.Text(), " ")

		if curr_instr[0] == "XX" { // if SPACE
			program = append(program, shared.Word(0))
			continue
		}

		binOp, err := strconv.Atoi(curr_instr[0])
		if err != nil {
			panic(err)
		}
		program = append(program, shared.Word(binOp))

		if len(curr_instr) > 1 {
			binOperand1, err := strconv.Atoi(curr_instr[1])
			if err != nil {
				panic(err)
			}

			program = append(program, shared.Word(binOperand1))
		}

		if len(curr_instr) > 2 {
			binOperand2, err := strconv.Atoi(curr_instr[2])
			if err != nil {
				panic(err)
			}
			program = append(program, shared.Word(binOperand2))
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	return program
}
