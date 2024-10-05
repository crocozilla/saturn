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

type SegmentSizes struct {
	text  []int
	data  []int
	space []int
}

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
	globalSymbolTable, segmentSizes :=
		firstPass(definitionTables, useTables, programNames, programSizes)

	fmt.Println("after:")
	fmt.Println(globalSymbolTable)
	for i := range programSizes {
		fmt.Println(useTables[i], programSizes[i])
	}

	// second pass here
	secondPass(useTables, programNames, globalSymbolTable, segmentSizes)

	totalStackSize := uint16(0)
	for _, size := range stackSizes {
		totalStackSize += size
	}
	return totalStackSize
}

func firstPass(
	definitionTables []map[string]shared.SymbolInfo,
	useTables []map[string][]uint16,
	programNames []string,
	programSizes []uint16) (
	globalSymbolTable map[string]shared.SymbolInfo,
	segmentSizes SegmentSizes) {

	globalSymbolTable = map[string]shared.SymbolInfo{}
	for i := range programSizes {
		programFile, err := shared.OpenBuildFile(programNames[i] + ".obj")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(programFile)
		textSize := 0
		dataSize := 0
		spaceSize := 0
		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			if lineFields[0] == "XX" {
				spaceSize++

				// only data and space can have size 2
			} else if len(lineFields) == 2 {
				dataSize++

			} else {
				for _, field := range lineFields {
					if field != "A" && field != "R" {
						textSize++
					}
				}
			}
		}

		segmentSizes.text = append(segmentSizes.text, textSize)
		segmentSizes.data = append(segmentSizes.data, dataSize)
		segmentSizes.space = append(segmentSizes.space, spaceSize)
	}

	for program_idx := range programSizes {
		// update useTables to global addresses
		useTable := useTables[program_idx]
		for symbol, uses := range useTable {
			for use := range uses {
				address := int(useTable[symbol][use])
				new_address := relocateRelativeAddress(
					address,
					program_idx,
					segmentSizes)
				useTable[symbol][use] = uint16(new_address)
			}
		}

		// copy definition tables to global symbol table
		// with correct global address
		definitionTable := definitionTables[program_idx]
		for symbol, info := range definitionTable {
			if _, ok := globalSymbolTable[symbol]; ok {
				panic(errors.New("símbolo global já definido"))
			}

			globalAddress := info.Address
			if info.Mode == shared.RELATIVE {
				address := int(globalAddress)
				new_address := relocateRelativeAddress(
					address,
					program_idx,
					segmentSizes)
				globalAddress = uint16(new_address)
			}

			globalSymbolTable[symbol] = shared.SymbolInfo{
				Address: globalAddress,
				Mode:    'A'}
		}
	}
	// check if all used symbols were defined
	for _, useTable := range useTables {
		for symbol := range useTable {
			if _, defined := globalSymbolTable[symbol]; !defined {
				panic(errors.New("simbolo " + symbol + " nao foi definido"))
			}
		}
	}

	return globalSymbolTable, segmentSizes
}

func secondPass(
	useTables []map[string][]uint16,
	programNames []string,
	globalSymbolTable map[string]shared.SymbolInfo,
	segmentSizes SegmentSizes) {

	hpxFile, err := shared.CreateBuildFile(programNames[0] + ".hpx")
	if err != nil {
		panic(err)
	}
	defer hpxFile.Close()

	locationCounter := 0

	sizeOfPreviousText := 0
	sizeOfPreviousData := 0

	// write text
	for program_idx, name := range programNames {
		programFile, err := shared.OpenBuildFile(name + ".obj")
		if err != nil {
			panic(err)
		}
		defer programFile.Close()
		scanner := bufio.NewScanner(programFile)
		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			updateLineFieldsAddresses(lineFields,
				globalSymbolTable,
				useTables[program_idx],
				&locationCounter,
				segmentSizes,
				program_idx)
			writeHpxLine(hpxFile, lineFields)
			if locationCounter >= sizeOfPreviousText+segmentSizes.text[program_idx] {
				break
			}
		}
		sizeOfPreviousText += segmentSizes.text[program_idx]
	}

	//write data
	for program_idx, name := range programNames {
		programFile, err := shared.OpenBuildFile(name + ".obj")
		if err != nil {
			panic(err)
		}
		defer programFile.Close()
		scanner := bufio.NewScanner(programFile)

		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			// skip text
			if len(lineFields) != 2 {
				continue
			}
			updateLineFieldsAddresses(lineFields,
				globalSymbolTable,
				useTables[program_idx],
				&locationCounter,
				segmentSizes,
				program_idx)
			writeHpxLine(hpxFile, lineFields)
			totalTextSize := sizeOfPreviousText
			if locationCounter >= totalTextSize+
				sizeOfPreviousData+segmentSizes.data[program_idx] {
				break
			}
		}
		sizeOfPreviousData += segmentSizes.data[program_idx]
	}

	// write space
	for program_idx, name := range programNames {
		programFile, err := shared.OpenBuildFile(name + ".obj")
		if err != nil {
			panic(err)
		}
		defer programFile.Close()
		programFile.Seek(0, segmentSizes.text[program_idx]+
			segmentSizes.data[program_idx])
		scanner := bufio.NewScanner(programFile)

		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			// skip text and data
			if lineFields[0] != "XX" {
				continue
			}
			updateLineFieldsAddresses(lineFields,
				globalSymbolTable,
				useTables[program_idx],
				&locationCounter,
				segmentSizes,
				program_idx)
			writeHpxLine(hpxFile, lineFields)
		}
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
	segmentSizes SegmentSizes,
	program_idx int) {

	for i := range lineFields {
		// 00 A is sentinel value for INTDEF/INTUSE? value
		if lineFields[i] == "00" {
			// cant break bounds because of way obj files are created
			if lineFields[i+1] == "A" {
				for symbol, useAddresses := range useTable {
					fmt.Println(symbol, useAddresses)
					for _, useAddress := range useAddresses {
						if useAddress == uint16(*locationCounter) {
							address := strconv.Itoa(
								int(globalSymbolTable[symbol].Address))

							lineFields[i] = address
						}
					}
				}
			}
		}
		if lineFields[i] == "R" {
			address, _ := strconv.Atoi(lineFields[i-1])
			new_address := relocateRelativeAddress(
				address,
				program_idx,
				segmentSizes)
			lineFields[i-1] = strconv.Itoa(new_address)

		}
		if lineFields[i] != "A" && lineFields[i] != "R" {
			(*locationCounter)++
		}
	}
}

func relocateRelativeAddress(
	address,
	program_idx int,
	segmentSizes SegmentSizes) int {

	textSize := segmentSizes.text[program_idx]
	dataSize := segmentSizes.data[program_idx]
	totalTextSize := 0
	totalDataSize := 0
	totalSpaceSize := 0
	sizeOfPreviousText := 0
	sizeOfPreviousData := 0
	sizeOfPreviousSpace := 0
	for i := range segmentSizes.text {
		totalTextSize += segmentSizes.text[i]
		totalDataSize += segmentSizes.data[i]
		totalSpaceSize += segmentSizes.space[i]
		if i < program_idx {
			sizeOfPreviousText += segmentSizes.text[i]
			sizeOfPreviousData += segmentSizes.data[i]
			sizeOfPreviousSpace += segmentSizes.space[i]
		}
	}

	isText := address < textSize
	isData := address >= textSize && address < textSize+dataSize
	otherTextSize := totalTextSize - textSize
	otherDataSize := totalDataSize - dataSize
	if isText {
		address += sizeOfPreviousText
	} else if isData {
		address += otherTextSize + sizeOfPreviousData
	} else {
		address += otherTextSize + otherDataSize + sizeOfPreviousSpace
	}
	return address
}
