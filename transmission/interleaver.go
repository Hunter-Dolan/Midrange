package transmission

import "math"

func interleave(data []bool) []bool {
	dataLength := len(data)

	rounds := int(math.Floor(float64(dataLength / 4)))

	for round := 0; round < rounds; round++ {
		shiftedData := make([]bool, dataLength)

		for index := range data {
			offset := int(math.Floor(float64(index / 2)))
			odd := index % 2

			if odd == 1 {
				shiftedData[index] = data[offset]
			} else {
				shiftedData[index] = data[dataLength-1-offset]
			}

		}

		data = shiftedData
	}

	return data
}
