package main

import (
	"fmt"
)

func main() {
	bitMask := 0b001100

	code := 0b000100

	fmt.Printf("%d", bitMask & code)
}
