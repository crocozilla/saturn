package main

import (
	"bufio"
	"os"
	"saturn/shared"
	"strconv"
	"strings"
)

func ReadProgram(filePath string) []shared.Word {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var program []shared.Word

	for scanner.Scan() {
		curr_instr := strings.Split(scanner.Text(), " ")

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
