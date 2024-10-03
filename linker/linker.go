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
		for symbol, info := range definitionTable {
			if _, ok := globalSymbolTable[symbol]; ok {
				panic(errors.New("símbolo global já definido"))
			}

			globalAddress := info.Address
			if info.Mode == shared.RELATIVE {
				globalAddress += sizeOfPreviousPrograms
			}

			globalSymbolTable[symbol] = shared.SymbolInfo{
				Address: globalAddress,
				Mode:    info.Mode}
		}
		sizeOfPreviousPrograms += programSizes[i]
	}
	// conferir se todos os simbolos usados foram definidos
	for _, useTable := range useTables {
		for symbol := range useTable {
			if _, defined := globalSymbolTable[symbol]; !defined {
				panic(errors.New("simbolo " + symbol + " nao foi definido"))
			}
		}
	}
	fmt.Println("after:")
	fmt.Println(globalSymbolTable)
	for i := range programSizes {
		fmt.Println(useTables[i], programSizes[i])
	}
}
