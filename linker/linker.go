package linker

import (
	"errors"
	"fmt"
	"os"
	"saturn/shared"
)

func Run(
	objFiles []*os.File,
	definitionTables []map[string]shared.SymbolInfo,
	useTables []map[string][]uint16,
	programSizes []uint16) {

	fmt.Println("")
	fmt.Println("before: ")
	for i := range programSizes {
		fmt.Println(definitionTables[i], useTables[i], programSizes[i])
	}
	globalSymbolTable := map[string]shared.SymbolInfo{}
	sizeOfPreviousPrograms := uint16(0)
	for i := range programSizes {
		// update useTables to global addresses
		useTable := useTables[i]
		for symbol, uses := range useTable {
			for use := range uses {
				useTable[symbol][use] += sizeOfPreviousPrograms
			}
		}

		// copy definition tables to global symbol table
		// with correct global address
		definitionTable := definitionTables[i]
		for key, value := range definitionTable {

			globalAddress := value.Address
			if value.Mode == shared.RELATIVE {
				globalAddress += sizeOfPreviousPrograms
			}
			if _, ok := globalSymbolTable[key]; ok {
				panic(errors.New("símbolo global já definido"))
			}
			globalSymbolTable[key] = value
		}
		sizeOfPreviousPrograms += programSizes[i]
	}

	fmt.Println("after:")
	fmt.Println(globalSymbolTable)
	for i := range programSizes {
		fmt.Println(useTables[i], programSizes[i])
	}
}
