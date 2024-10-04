package mp

// Macro Processor
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"saturn/parser"
	"saturn/shared"
	"slices"
)

type macroInstructions []string
type macro struct {
	numberOfParameters int
	instructions       macroInstructions
}

type macroProcessor struct {
	macroDefinitiontable map[string]macro
}

func New() *macroProcessor {
	macroProcessor := new(macroProcessor)
	macroProcessor.macroDefinitiontable = map[string]macro{}
	return macroProcessor
}

func (macroProcessor *macroProcessor) MacroPass(file *os.File) *os.File {
	scanner := bufio.NewScanner(file)

	masmaprg, err := shared.CreateBuildFile("MASMAPRG.ASM")
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		line, isComment := parser.ReadLine(scanner)
		if isComment {
			continue
		}

		label, operationString, operands := parser.MacroLine(line)

		if operationString == "MACRO" {
			macroProcessor.macroDefine(scanner)
			continue

		} else if _, ok := macroProcessor.macroDefinitiontable[operationString]; ok {
			macroProcessor.macroExpand(line, masmaprg)
			continue
		}

		// write line to file:
		var operandsString string
		for i := range operands {
			operandsString += operands[i] + " "
		}
		writtenLine := fmt.Sprintln(label, operationString, operandsString)
		masmaprg.WriteString(writtenLine)
	}

	return masmaprg
}

// gets the macro definition, starts after "MACRO" operation/instruction
func (macroProcessor *macroProcessor) macroDefine(scanner *bufio.Scanner) {

	var macro macro
	var macroName string
	var macroOperands []string
	isDefinition := true // first line after MACRO
	quit := false
	definitionLevel := 1
	isFirstDefinition := true
	parameterStack := [][2]string{}
	for scanner.Scan() && !quit {

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

			if isFirstDefinition {
				macroName = currentName
				macro.numberOfParameters = len(macroOperands)
				isFirstDefinition = false
				continue
			}

		}

		label, operationString, lineOperands := parser.MacroLine(line)
		if operationString == "MACRO" {
			definitionLevel++
			isDefinition = true

		} else if operationString == "MEND" {
			parameterStack = popLevelFromStack(parameterStack)
			definitionLevel--
			if definitionLevel == 0 {
				quit = true
			}
		}

		replaceNamesByCodes(parameterStack, &label, &operationString, lineOperands)

		macroLine := createMacroLine(label, operationString, lineOperands)
		macro.instructions = append(macro.instructions, macroLine)

	}

	if quit {
		macroProcessor.macroDefinitiontable[macroName] = macro
	} else {
		panic("faltando diretiva MEND")
	}

}

func (macroProcessor *macroProcessor) macroDefineFromSlice(
	macroInstructions []string,
	parameterStack [][2]string) int {

	var macro macro
	var macroName string
	var macroOperands []string
	isDefinition := true // first line after MACRO
	isFirstDefinition := true
	quit := false
	initialDefinitionLevel := 2 // this function is only called when level is already 1
	definitionLevel := initialDefinitionLevel
	idx := 0
	var line string
	for idx, line = range macroInstructions {

		if quit {
			break
		}

		if isDefinition {
			var currentName string
			var operand0 string
			operand0, currentName, macroOperands = parser.MacroLine(line)
			if operand0 != "" {
				macroOperands = slices.Insert(macroOperands, 0, operand0)
			}
			// doesnt check if parameters start with & because our internal representation
			// doesnt use &

			parameterStack = addToStack(parameterStack, definitionLevel, macroOperands)

			isDefinition = false

			if isFirstDefinition {
				macroName = currentName
				macro.numberOfParameters = len(macroOperands)
				isFirstDefinition = false
				continue
			}

		}

		label, operationString, lineOperands := parser.MacroLine(line)
		if operationString == "MACRO" {
			definitionLevel++
			isDefinition = true

		} else if operationString == "MEND" {
			parameterStack = popLevelFromStack(parameterStack)
			definitionLevel--
			if definitionLevel == initialDefinitionLevel-1 {
				quit = true
			}
		}

		replaceCodesByNames(parameterStack, &label, &operationString, lineOperands)

		macroLine := createMacroLine(label, operationString, lineOperands)
		macro.instructions = append(macro.instructions, macroLine)

	}

	if quit {
		macroProcessor.macroDefinitiontable[macroName] = macro
		return idx
	} else {
		panic("faltando diretiva MEND")
	}

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
	for idx := 0; idx < len(macro.instructions); idx++ {
		instructionLine := macro.instructions[idx]
		label, operationString, operands := parser.MacroLine(instructionLine)
		replaceCodesByNames(parameterStack, &label, &operationString, operands)
		removeAmpersands(&label, &operationString, operands)

		macroLine := createMacroLine(label, operationString, operands)

		if _, isMacro := macroProcessor.macroDefinitiontable[operationString]; isMacro {
			macroProcessor.macroExpand(macroLine, masmaprg)
			continue
		}
		if operationString == "MACRO" {
			idx += macroProcessor.macroDefineFromSlice(macro.instructions[idx+1:], parameterStack)
			continue
		}
		if operationString == "MEND" {
			return
		}

		masmaprg.WriteString(macroLine + "\n")
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

func matchInStack(
	parameterStack [][2]string,
	token string,
	tokensAreNames bool) (replacement string, valid bool) {

	name := 0
	code := 1
	tokensAreCodes := !tokensAreNames

	// for start at the end to get closest scope
	if tokensAreNames {
		for i := len(parameterStack) - 1; i >= 0; i-- {
			if parameterStack[i][name] == token {
				return parameterStack[i][code], true
			}
		}
	} else if tokensAreCodes {
		for i := len(parameterStack) - 1; i >= 0; i-- {
			if parameterStack[i][code] == token {
				return parameterStack[i][name], true
			}
		}
	}

	return "", false
}

func addToStack(
	parameterStack [][2]string,
	definitionLevel int,
	operands []string) (NewParameterStack [][2]string) {

	code := "#(1,1)"
	codeLength := len(code)
	for i, op := range operands {
		// replaces temporary #(2,1) or #(3,1) by permanent #(1,1), if its the case
		if len(op) == codeLength {
			if op[:4] == fmt.Sprintf("#(%d,", definitionLevel) {
				op = fmt.Sprintf("#(1,%v)", getDigit(op[4]))
			}
		}
		//

		parameter := [2]string{op, fmt.Sprintf("#(%d,%d)", definitionLevel, i+1)}
		parameterStack = append(parameterStack, parameter)
	}
	return parameterStack
}

// replaces #(1,1) #(1,2) for ARG1 ARG2
func replaceCodesByNames(parameterStack [][2]string, label, operation *string, operands []string) {
	tokensAreNames := false
	replaceTokens(parameterStack, label, operation, operands, tokensAreNames)
}

// replaces ARG1 ARG2 for #(1,1) #(1,2)
func replaceNamesByCodes(parameterStack [][2]string, label, operation *string, operands []string) {
	tokensAreNames := true
	replaceTokens(parameterStack, label, operation, operands, tokensAreNames)
}

// tokens are names tells stack if we should sub ARG1 for #1 or vice versa
func replaceTokens(
	parameterStack [][2]string,
	label,
	operation *string,
	operands []string, tokensAreNames bool) {

	if replacement, valid := matchInStack(parameterStack, *label, tokensAreNames); valid {
		*label = fmt.Sprintf("%v", replacement)
	}
	if replacement, valid := matchInStack(parameterStack, *operation, tokensAreNames); valid {
		*operation = fmt.Sprintf("%v", replacement)
	}

	for i := range operands {
		if replacement, valid := matchInStack(parameterStack, operands[i], tokensAreNames); valid {
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

func getDigit(char byte) int {
	return int(char) - '0'
}

// deletes parameters from latest scope
func popLevelFromStack(parameterStack [][2]string) [][2]string {
	code := 1
	levelInString := 2
	var definitionLevel int
	if len(parameterStack) > 0 {
		lastParameter := parameterStack[len(parameterStack)-1]
		definitionLevel = getDigit(lastParameter[code][levelInString])
	}
	for i := len(parameterStack) - 1; i >= 0; i-- {
		digit := getDigit(parameterStack[i][code][levelInString])
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
