package mp

// Macro Processor
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"saturn/parser"
)

/*
 O processamento de macros deve ser realizado antes da montagem, sendo ativado a partir do
 módulo principal integrador do macro-montador.

 Deve permitir a definição de macro dentro de macros (macros aninhadas),
 bem como a chamada de macros dentro de macro (chamadas aninhadas), sendo,
 portanto, implementado em uma só passagem. O programa
 receberá como entrada um arquivo fonte informado para montagem e gerará como saída
 outro arquivo fonte com o nome  MASMAPRG.ASM.

 As macros são definidas através das pseudo-operações MACRO e MEND e a sintaxe
 está exemplificada no Anexo 2 por um programa que deverá ser
 utilizado para teste do Processador de Macros.

 Também deve ser prevista a opção de listagem das macros segundo formato a ser
 combinado com o professor. Esta listagem deverá também conter algumas estatísticas sobre
 o uso das macros.
*/

type macroInstructions []string
type macro struct {
	numberOfParameters int
	instructions       macroInstructions
}

type macroProcessor struct {
	lineCounter          uint16
	macroDefinitiontable map[string]macro
}

func New() *macroProcessor {
	macroProcessor := new(macroProcessor)
	macroProcessor.macroDefinitiontable = map[string]macro{}
	return macroProcessor
}

func (macroProcessor *macroProcessor) MacroPass(file *os.File) {
	scanner := bufio.NewScanner(file)

	masmaprg, err := os.Create("MASMAPRG.ASM")
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		macroProcessor.lineCounter++
		line, isComment := parser.ReadLine(scanner)
		if isComment {
			continue
		}

		label, operationString, operands := parser.MacroLine(line)
		var operandsString string

		if operationString == "MACRO" {
			macroProcessor.macroDefine(scanner)
			continue
		}

		// write line to file:
		for i := range operands {
			operandsString += operands[i]
			operandsString += " "
		}
		writtenLine := fmt.Sprintln(label, operationString, operandsString)
		masmaprg.WriteString(writtenLine)
	}
}

// gets the macro definition, starts after "MACRO" operation/instruction
func (macroProcessor *macroProcessor) macroDefine(scanner *bufio.Scanner) {

	var macro macro
	var macroName string
	var macroOperands []string
	isDefinition := true // first line after MACRO
	definitionLevel := 1
	for scanner.Scan() {
		macroProcessor.lineCounter++

		line, isComment := parser.ReadLine(scanner)
		if isComment {
			continue
		}

		if isDefinition {
			_, macroName, macroOperands = parser.MacroLine(line)
			// what to do with label?

			macro.numberOfParameters = len(macroOperands)

			err := checkMacroOperands(macroOperands)
			if err != nil {
				panic(err)
			}
			isDefinition = false
			continue
		}

		label, operationString, lineOperands := parser.MacroLine(line)
		label, operationString, lineOperands = substituteOperands(label, operationString, lineOperands, macroOperands)

		if operationString == "MACRO" {
			definitionLevel++

		} else if operationString == "MEND" {
			definitionLevel--
			if definitionLevel == 0 {
				macroProcessor.macroDefinitiontable[macroName] = macro
				return
			}
		}

		macroLine := createMacroLine(label, operationString, lineOperands)
		macro.instructions = append(macro.instructions, macroLine)

	}

	panic("faltando diretiva MEND")

}

func createMacroLine(label, operation string, operands []string) string {
	var macroLine string
	macroLine += label + " "
	macroLine += operation + " "
	for _, op := range operands {
		macroLine += op + " "
	}

	return macroLine
}

func checkMacroOperands(operands []string) error {
	for _, op := range operands {
		if op[0] != '&' {
			return (errors.New("operandos para macros devem começar com '&'"))
		}
	}
	return nil
}

func substituteOperands(label, operation string, lineOperands, macroOperands []string) (string, string, []string) {
	for i, op := range macroOperands {
		if label == op {
			label = fmt.Sprintf("#%d", i+1)
		}
		if operation == op {
			operation = fmt.Sprintf("#%d", i+1)
		}

		for i := range lineOperands {
			if lineOperands[i] == op {
				lineOperands[i] = fmt.Sprintf("#%d", i+1)
			}
		}
	}

	return label, operation, lineOperands

}

func (macroProcessor *macroProcessor) macroExpand(name string) {

}
