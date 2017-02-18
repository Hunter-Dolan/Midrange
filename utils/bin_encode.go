package utils

import "fmt"

func BinEncode(s []byte) []bool {
	binArray := make([]bool, len(s)*8)

	binIndex := 0

	for _, c := range s {
		binary := fmt.Sprintf("%b", c)
		binaryLen := len(binary)

		padAmount := 8 - binaryLen

		for i := 0; i < padAmount; i++ {
			binArray[binIndex] = false
			binIndex++
		}

		for _, bit := range binary {
			binArray[binIndex] = bit == '1'
			binIndex++
		}
	}

	return binArray
}

func BinEncodeWithPad(s []byte, pad int) []bool {
	bin := BinEncode(s)

	paddedBin := make([]bool, pad-len(bin))
	paddedBin = append(paddedBin, bin...)

	return paddedBin
}

func BoolBinEncode(s string) []bool {
	bin := make([]bool, len(s))

	for i, value := range s {
		bin[i] = value == '1'
	}

	return bin
}
