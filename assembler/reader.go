package assembler

import (
	"bufio"
	"os"
)

func scanWords(path string, callback func(word string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		callback(word)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
