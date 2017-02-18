package fec

import "math"

func (d *DataConfidenceCollection) Deinterleave() {

	confidence := d.confidence
	data := d.data

	dataLength := len(data)
	dataOdd := dataLength % 2

	rounds := int(math.Floor(float64(dataLength / 4)))

	for round := 0; round < rounds; round++ {
		shiftedData := make([]bool, dataLength)
		shiftedConfidence := make([]float64, dataLength)

		for index := range data {
			pastMidpoint := (float64(index+1) / float64(dataLength)) > float64(0.5)

			var shiftedIndex int

			if pastMidpoint {
				shiftedIndex = (dataLength - (index+1-int(math.Floor(float64(dataLength)/2.0)))*2)
				if dataOdd == 1 {
					shiftedIndex++
				}
			} else {
				shiftedIndex = index*2 + 1
			}

			shiftedData[index] = data[shiftedIndex]
			shiftedConfidence[index] = confidence[shiftedIndex]
		}

		data = shiftedData
		confidence = shiftedConfidence
	}

	d.data = data
	d.confidence = confidence
}
