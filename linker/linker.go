package linker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"saturn/shared"
	"strconv"
	"strings"
)

func Run(
	definitionTables []map[string]shared.SymbolInfo,
	useTables []map[string][]uint16,
	programNames []string,
	programSizes []uint16,
	stackSizes []uint16) uint16 {

	if len(definitionTables) == 0 {
		return 0
	}

	fmt.Println("")
	fmt.Println("before: ")
	for i := range programSizes {
		fmt.Println(definitionTables[i], useTables[i], programSizes[i])
	}
	globalSymbolTable := firstPass(definitionTables, useTables, programSizes)

	fmt.Println("after:")
	fmt.Println(globalSymbolTable)
	for i := range programSizes {
		fmt.Println(useTables[i], programSizes[i])
	}

	// second pass here
	secondPass(useTables, programNames, programSizes, globalSymbolTable)

	totalStackSize := uint16(0)
	for _, size := range stackSizes {
		totalStackSize += size
	}
	return totalStackSize
}

func firstPass(
	definitionTables []map[string]shared.SymbolInfo,
	useTables []map[string][]uint16,
	programSizes []uint16) map[string]shared.SymbolInfo {

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

	return globalSymbolTable
}

func secondPass(
	useTables []map[string][]uint16,
	programNames []string,
	programSizes []uint16,
	globalSymbolTable map[string]shared.SymbolInfo) {

	hpxFile, err := shared.CreateBuildFile(programNames[0] + ".hpx")
	if err != nil {
		panic(err)
	}
	defer hpxFile.Close()

	var scanner *bufio.Scanner
	locationCounter := 0
	sizeOfPreviousPrograms := uint16(0)
	for program_idx, name := range programNames {
		programFile, err := shared.OpenBuildFile(name + ".obj")
		if err != nil {
			panic(err)
		}
		scanner = bufio.NewScanner(programFile)
		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			updateLineFieldsAddresses(lineFields,
				globalSymbolTable,
				useTables[program_idx],
				&locationCounter,
				sizeOfPreviousPrograms)
			writeHpxLine(hpxFile, lineFields)
		}
		sizeOfPreviousPrograms += programSizes[program_idx]

	}
}

func writeHpxLine(hpxFile *os.File, lineFields []string) {
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

// updates external addresses (00 A) to actual addresses
// updates locationCounter
func updateLineFieldsAddresses(
	lineFields []string,
	globalSymbolTable map[string]shared.SymbolInfo,
	useTable map[string][]uint16,
	locationCounter *int,
	sizeOfPreviousPrograms uint16) {
	for i := range lineFields {
		// 00 A is sentinel value for INTDEF/INTUSE? value
		if lineFields[i] == "00" {
			// cant break bounds because of way obj files are created
			if lineFields[i+1] == "A" {
				for symbol, useAddresses := range useTable {
					for _, useAddress := range useAddresses {
						if useAddress == uint16(*locationCounter) {
							address := strconv.Itoa(
								int(globalSymbolTable[symbol].Address))
							mode := string(globalSymbolTable[symbol].Mode)

							lineFields[i] = address
							lineFields[i+1] = mode
						}
					}
				}
			}
		}
		if lineFields[i] == "R" {
			fieldValue, _ := strconv.Atoi(lineFields[i-1])
			fieldValue += int(sizeOfPreviousPrograms)
			lineFields[i-1] = strconv.Itoa(fieldValue)
		}
		if lineFields[i] != "A" && lineFields[i] != "R" {
			(*locationCounter)++
		}
	}
}
