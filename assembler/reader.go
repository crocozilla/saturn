package assembler

import (
	"bufio"
	"os"
)

// TODO: Trocar para ler por linha
func scanLines(path string, callback func(line string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		callback(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
