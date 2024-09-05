package assembler

import (
	"math"
	"saturn/shared"
	"strconv"
	"testing"
)

/*
func TestFirstPass(t *testing.T) {
	file, err := os.Open("first_pass_test.asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	assembler := New()

	assembler.firstPass(file)
}*/

// missing literal tests
func TestGetOperandValue(t *testing.T) {
	iTest, err := getOperandValue("#64")
	if err != nil {
		panic(err)
	}
	if iTest != 64 {
		t.Fatalf("número imediato inválido")
	}
	hexTest, err := getOperandValue("H'F'")
	if err != nil {
		panic(err)
	}
	if hexTest != 15 {
		t.Fatalf("número hexadecimal inválido, esperava-se %v mas obteve-se %v", 15, hexTest)
	}
	hexTest, err = getOperandValue("H'32AF'")
	if err != nil {
		panic(err)
	}
	if hexTest != 12975 {
		t.Fatalf("número hexadecimal inválido, esperava-se %v mas obteve-se %v", 12975, hexTest)
	}
	//test, err = getOperandValue("@30")
	maxNumber := int(math.Pow(2.0, float64(shared.WordSize-1)) - 1.0)
	minNumber := int(-math.Pow(2.0, float64(shared.WordSize-1)))

	_, err = getOperandValue(strconv.Itoa(maxNumber + 1))
	if err == nil {
		t.Fatalf("número grande demais não gerou erro")
	}
	_, err = getOperandValue(strconv.Itoa(maxNumber))
	if err != nil {
		t.Fatalf("número máximo gerou erro incorretamente")
	}
	_, err = getOperandValue(strconv.Itoa(minNumber - 1))
	if err == nil {
		t.Fatalf("número pequeno demais não gerou erro")
	}
	_, err = getOperandValue(strconv.Itoa(minNumber))
	if err != nil {
		t.Fatalf("número mínimo negativo gerou erro incorretamente")
	}
	_, err = getOperandValue("-")
	if err == nil {
		t.Fatalf("número negativo vazio não gerou erro")
	}
}
