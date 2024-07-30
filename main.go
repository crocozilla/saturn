package main

import (
	"fmt"
	"saturn/shared"
)

func main() {
	i := shared.Instruction{ 
		Operation: shared.ADD,
		Operands: shared.Operands{
			First: 2,
			Second: 5,
		},
	}

	fmt.Println(i)
}
