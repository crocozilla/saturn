package main

//TODO: código a ser executado deve ser guardado na maquina
//TODO: leitura de arquivo txt para pegar o codigo

//"saturn/gui"
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
	var instructions shared.Word

	for scanner.Scan() {
		curr_instr := strings.Split(scanner.Text(), " ")

		binOp, err := strconv.Atoi(curr_instr[0])
		if err != nil {
			panic(err)
		}
		binOperand1, err := strconv.Atoi(curr_instr[1])
		if err != nil {
			panic(err)
		}
		binOperand2, err := strconv.Atoi(curr_instr[2])
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, shared.Word(binOp))
		instructions = append(instructions, shared.Word(binOperand1))
		instructions = append(instructions, shared.Word(binOperand2))

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	return instructions
}

func main() {
	ReadProgram("./program.txt")
	//gui.Run()
}
