package mp

// Macro Processor
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"saturn/parser"
	"slices"
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

		} else if _, ok := macroProcessor.macroDefinitiontable[operationString]; ok {
			macroProcessor.macroExpand(line, masmaprg)
			continue
		}

		// write line to file:
		for i := range operands {
			operandsString += operands[i] + " "
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
	quit := false
	definitionLevel := 1
	parameterStack := [][2]string{}
	for scanner.Scan() && !quit {
		macroProcessor.lineCounter++

		line, isComment := parser.ReadLine(scanner)
		if isComment {
			continue
		}

		if isDefinition {
			var currentName string
			var operand0 string
			operand0, currentName, macroOperands = parser.MacroLine(line)
			if operand0 != "" {
				macroOperands = slices.Insert(macroOperands, 0, operand0)
			}
			err := checkMacroOperands(macroOperands)
			if err != nil {
				panic(err)
			}
			parameterStack = addToStack(parameterStack, definitionLevel, macroOperands)

			isDefinition = false

			if definitionLevel == 1 {
				macroName = currentName
				macro.numberOfParameters = len(macroOperands)
				continue
			}

		}

		label, operationString, lineOperands := parser.MacroLine(line)
		if operationString == "MACRO" {
			definitionLevel++
			isDefinition = true

		} else if operationString == "MEND" {
			parameterStack = deleteLevelFromStack(parameterStack, definitionLevel)
			definitionLevel--
			if definitionLevel == 0 {
				quit = true
			}
		}

		replaceTokens(parameterStack, &label, &operationString, lineOperands)

		macroLine := createMacroLine(label, operationString, lineOperands)
		macro.instructions = append(macro.instructions, macroLine)

	}

	if quit {
		macroProcessor.macroDefinitiontable[macroName] = macro
	} else {
		panic("faltando diretiva MEND")
	}

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

func matchInStack(parameterStack [][2]string, token string) (replacement string, valid bool) {
	name := 0
	code := 1
	// starts at the end to get closest scope
	for i := len(parameterStack) - 1; i >= 0; i-- {
		if parameterStack[i][name] == token {
			return parameterStack[i][code], true

		} else if parameterStack[i][code] == token {
			return parameterStack[i][name], true
		}
	}

	return "", false
}

func addToStack(parameterStack [][2]string, definitionLevel int, operands []string) (NewParameterStack [][2]string) {
	for i, op := range operands {
		parameterStack = append(parameterStack, [2]string{op, fmt.Sprintf("#(%d,%d)", definitionLevel, i+1)})
	}
	return parameterStack
}

func replaceTokens(parameterStack [][2]string, label, operation *string, operands []string) {
	if replacement, valid := matchInStack(parameterStack, *label); valid {
		*label = fmt.Sprintf("%v", replacement)
	}
	if replacement, valid := matchInStack(parameterStack, *operation); valid {
		*operation = fmt.Sprintf("%v", replacement)
	}

	for i := range operands {
		if replacement, valid := matchInStack(parameterStack, operands[i]); valid {
			operands[i] = fmt.Sprintf("%v", replacement)
		}
	}
}

func removeAmpersands(label, operation *string, operands []string) {
	if len(*label) > 0 {
		if (*label)[0] == '&' {
			*label = (*label)[1:]
		}
	}
	if (*operation)[0] == '&' {
		*operation = (*operation)[1:]
	}
	for i := range operands {
		if operands[i][0] == '&' {
			operands[i] = operands[i][1:]
		}
	}
}

func deleteLevelFromStack(parameterStack [][2]string, definitionLevel int) [][2]string {
	level := 2

	for i := len(parameterStack) - 1; i >= 0; i-- {
		digit := int(parameterStack[i][1][level]) - '0'
		fmt.Println(digit)
		if digit == definitionLevel {
			parameterStack = parameterStack[:len(parameterStack)-1]

		} else {
			return parameterStack
		}
	}

	return [][2]string{}
}

func checkMacroOperands(operands []string) error {
	for _, op := range operands {
		if op[0] != '&' {
			return (errors.New("operandos para macros devem começar com '&'"))
		}
	}
	return nil
}

func (macroProcessor *macroProcessor) macroExpand(line string, masmaprg *os.File) {
	operand0, name, operands := parser.MacroLine(line)
	macro := macroProcessor.macroDefinitiontable[name]
	if operand0 != "" {
		operands = slices.Insert(operands, 0, operand0)
	}
	if len(operands) > macro.numberOfParameters {
		panic("um macro tem parametros demais")
	}

	parameterStack := [][2]string{}
	parameterStack = addToStack(parameterStack, 1, operands)

	// substitutes things like #1 #2 for arg1 arg2
	for _, instructionLine := range macro.instructions {
		label, operationString, operands := parser.MacroLine(instructionLine)
		replaceTokens(parameterStack, &label, &operationString, operands)
		removeAmpersands(&label, &operationString, operands)

		macroLine := createMacroLine(label, operationString, operands)

		if _, isMacro := macroProcessor.macroDefinitiontable[operationString]; isMacro {
			macroProcessor.macroExpand(macroLine, masmaprg)
			continue
		}
		if operationString == "MEND" {
			return
		}

		masmaprg.WriteString(macroLine + "\n")
	}
}
