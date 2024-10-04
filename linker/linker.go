package linker

import (
	"bufio"
	"errors"
	"fmt"
	"saturn/shared"
	"strconv"
	"strings"
)

func Run(
	definitionTables []map[string]shared.SymbolInfo,
	useTables []map[string][]uint16,
	programNames []string,
	programSizes []uint16) {

	if len(definitionTables) == 0 {
		return
	}

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
	// check if all used symbols were defined
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

	hpxFile, err := shared.CreateBuildFile(programNames[0] + ".hpx")
	if err != nil {
		panic(err)
	}
	defer hpxFile.Close()

	var scanner *bufio.Scanner
	locationCounter := 0
	for program_idx, name := range programNames {
		programFile, err := shared.OpenBuildFile(name + ".obj")
		if err != nil {
			panic(err)
		}
		scanner = bufio.NewScanner(programFile)
		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			var line string

			for i := range lineFields {
				if lineFields[i] == "00" {
					// cant break bounds because of way obj files are created
					if lineFields[i+1] == "A" {
						for symbol, uses := range useTables[program_idx] {
							for _, use := range uses {
								if use == uint16(locationCounter) {
									address := globalSymbolTable[symbol].Address
									mode := globalSymbolTable[symbol].Mode
									addressString := strconv.Itoa(int(address))
									lineFields[i] = addressString
									lineFields[i+1] = string(mode)
								}
							}
						}
					}
				}
				line += lineFields[i]
				if lineFields[i] != "A" && lineFields[i] != "R" {
					locationCounter++
				}
			}
			var hpxLine string
			for i := range lineFields {
				if lineFields[i] != "A" && lineFields[i] != "R" {
					hpxLine += lineFields[i]
				}
				if i+1 < len(lineFields) {
					hpxLine += " "
				} else {
					hpxLine += "\n"
				}
			}
			hpxFile.WriteString(hpxLine)
		}

	}
}
