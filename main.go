package main

import (
	"fmt"
	"saturn/gui"
)

//TODO: código a ser executado deve ser guardado na maquina
//TODO: leitura de arquivo txt para pegar o codigo

//"saturn/gui"

func main() {
	program := ReadProgram("./program.txt")
	gui.InsertProgram(program)
	gui.Run()
}
