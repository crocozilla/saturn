package assembler

import (
	"unicode"
)

// TODO: Trocar para ler por linha
/*
func scanLines(assembler Assembler, path string, callback func(assembler Assembler, line string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		callback(assembler, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
*/

// starting at a non-space character, returns string up until a space
func getWord(line string) string {
	for i, v := range line {
		if unicode.IsSpace(v) {
			return line[:i]
		}
	}

	// returns here if word is at end of line, or if line is empty
	return line
}

// returns empty string if no words left
func skipUntilNextWord(line string) string {
	// skips current word
	line = line[len(getWord(line)):]

	for i, v := range line {
		if !unicode.IsSpace(v) {
			return line[i:]
		}
	}

	return ""
}

func beginsComment(line string) bool {
	if len(line) == 0 {
		return false
	}
	return line[0] == '*'
}

// assumes line is not a comment, if something optional is missing, returns empty string instead
func parseLine(line string) (label string, operation string, op1 string, op2 string) {
	label = getWord(line)

	line = skipUntilNextWord(line)
	operation = getWord(line)

	line = skipUntilNextWord(line)
	if beginsComment(line) {
		return label, operation, "", ""
	}
	op1 = getWord(line)

	line = skipUntilNextWord(line)
	if beginsComment(line) {
		return label, operation, op1, ""
	}
	op2 = getWord(line)

	line = skipUntilNextWord(line)
	if len(line) != 0 && !beginsComment(line) {
		panic("alguma linha tem colunas demais")
	}

	return label, operation, op1, op2
}
