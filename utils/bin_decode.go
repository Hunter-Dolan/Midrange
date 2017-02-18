package utils

import "strconv"

func BinDecode(binary []bool) []byte {
	bitLength := 8

	byteCount := len(binary) / bitLength
	byteArray := make([]byte, byteCount)

	for byteIndex := 0; byteIndex < byteCount; byteIndex++ {
		offset := byteIndex * bitLength
		boolBits := binary[offset : offset+bitLength]

		bits := ""

		for _, boolBit := range boolBits {
			if boolBit {
				bits += "1"
			} else {
				bits += "0"
			}
		}

		i, _ := strconv.ParseInt(bits, 2, 64)
		byteArray[byteIndex] = byte(i)
	}

	return byteArray
}
