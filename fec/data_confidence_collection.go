package fec

import "sort"

type DataConfidenceCollection struct {
	data       []bool
	confidence []float64

	redundantCollections []*DataConfidenceCollection
	maxConfidence        float64

	HashVerifier func(data []bool, hash []bool) bool
}

func (d *DataConfidenceCollection) Add(bit bool, bitConfidence float64) {
	d.data = append(d.data, bit)
	d.confidence = append(d.confidence, bitConfidence)
}

func (d DataConfidenceCollection) Data() []bool {
	return d.data
}

func (d DataConfidenceCollection) Get(index int) (bool, float64) {
	boolBit := d.data[index]
	confidence := d.confidence[index]

	return boolBit, confidence
}

func (d DataConfidenceCollection) Length() int {
	return len(d.data)
}

func (d DataConfidenceCollection) LeastConfidentBitIndexes(number int) []int {
	sortedConfidences := make([]float64, len(d.confidence))

	for i, confidence := range d.confidence {
		sortedConfidences[i] = confidence
	}

	sort.Float64s(sortedConfidences)

	returnData := []int{}

	for _, confidence := range sortedConfidences {
		index := -1

		for i, unsortedConfidence := range d.confidence {
			if unsortedConfidence == confidence {
				index = i
			}
		}

		returnData = append(returnData, index)

		if len(returnData) >= number {
			break
		}
	}

	return returnData
}

func (d DataConfidenceCollection) copy() *DataConfidenceCollection {
	return &d
}

func (d *DataConfidenceCollection) Append(appendage *DataConfidenceCollection) {
	d.data = append(d.data, appendage.data...)
	d.confidence = append(d.confidence, appendage.confidence...)
}

func (d *DataConfidenceCollection) Slice(leftEnd, length int) *DataConfidenceCollection {
	slice := DataConfidenceCollection{}

	rightEnd := int(leftEnd + length)
	dataLength := len(d.data)

	if rightEnd > dataLength {
		rightEnd = dataLength - 1
	}

	slice.data = d.data[leftEnd:rightEnd]
	slice.confidence = d.confidence[leftEnd:rightEnd]

	return &slice
}
