package fec

import "math"

func (d *DataConfidenceCollection) AddRedundantCollection(redundant *DataConfidenceCollection) {
	d.redundantCollections = append(d.redundantCollections, redundant)
}

func (d *DataConfidenceCollection) CorrectDataWithRedundancy() {
	for i, value := range d.data {
		zeroConfidence := 0.0
		oneConfidence := 0.0

		zeroCount := 0
		oneCount := 0

		if value {
			oneConfidence += d.confidence[i]
			oneCount++
		} else {
			zeroConfidence += d.confidence[i]
			zeroCount++
		}

		for _, redundantCollection := range d.redundantCollections {
			value := redundantCollection.data[i]

			if value {
				oneConfidence += d.confidence[i]
				oneCount++
			} else {
				zeroConfidence += d.confidence[i]
				zeroCount++
			}

		}

		if oneCount > zeroCount {
			d.data[i] = true
		} else {
			d.data[i] = false
		}

		d.confidence[i] = math.Abs(oneConfidence - zeroConfidence)
	}
}
