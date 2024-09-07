package main

import (
	"saturn/assembler"
	// "saturn/gui"
)

//TODO: código a ser executado deve ser guardado na maquina
//TODO: leitura de arquivo txt para pegar o codigo

//"saturn/gui"

func main() {
	// program := ReadProgram("./program.txt")
	// gui.InsertProgram(program)
	// gui.Run()

	assembler.Run("assembler/first_pass_test.asm")
}
