package assembler

import(
	"os"
	"bufio"
	"testing"
	"fmt"
)

func TestParseLine(t *testing.T){
	file, err := os.Open("parse_line_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan(){
		line := scanner.Text()
		label, operation, op1, op2 := parseLine(line)
		fmt.Println("line", i, ":")
		if(label == ""){
			label = "empty"
		}
		if(operation == ""){
			operation = "empty"
		}
		if(op1 == ""){
			op1 = "empty"
		}
		if(op2 == ""){
			op2 = "empty"
		}
		fmt.Println("  ", label, operation, op1, op2)
		i++
	}

}