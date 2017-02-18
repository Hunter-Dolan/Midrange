package utils

func LeftPad(str string, length int) string {
	return LeftPadWithByte(str, length, "0")
}

func LeftPadWithByte(str string, length int, paddingByte string) string {
	padding := ""
	paddingAmount := length - len(str)

	for i := 0; i < paddingAmount; i++ {
		padding += paddingByte
	}

	return padding + str
}
