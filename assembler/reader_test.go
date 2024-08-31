package assembler

import(
	"os"
	"bufio"
	"testing"
	"fmt"
)

func TestParseLine(t *testing.T){
	file, err := os.Open("reader_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := scanner.Text()
		label, operation, op1, op2 := parseLine(line)
		fmt.Println(label, operation, op1, op2)
	}

}