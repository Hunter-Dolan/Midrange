package utils

func BoolToBin(boolean []bool) string {
	result := ""

	for _, value := range boolean {
		if value {
			result += "1"
		} else {
			result += "0"
		}
	}

	return result
}

func BoolToByteBin(boolean []bool) []byte {
	result := make([]byte, len(boolean))

	for i, value := range boolean {
		if value {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}

	return result
}
