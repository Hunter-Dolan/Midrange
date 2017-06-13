package fec

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

func (d *DataConfidenceCollection) CorrectData(hash *DataConfidenceCollection, hashValidated bool) bool {

	// Check to see if the data is valid
	if d.verifyHash(hash) {
		return true
	}

	maxRounds := options.MaxCorrectionRoundsHashValidated

	if !hashValidated {
		maxRounds = options.MaxCorrectionRoundsHashNotValidated
	}

	originalHash := hash.copy()

	valid := false

	for round := 0; round < maxRounds; round++ {

		if round != 0 {
			fmt.Println("Corrector Round:", round)
		}

		bits := possibleBinarySlicesWithNumberOfBits(round + 1)
		bitsLen := len(bits)

		//par.ForInterleaved(0, uint(bitsLen-1), 1, func(index uint) {
		for index := 0; index < bitsLen; index++ {
			if !valid {
				localValid := false

				data := d

				dataBits := bits[index]

				// First Attempt to flip only data bits
				data.swapLeastConfidentBitsWithBits(dataBits)

				if data.verifyHash(originalHash) {
					localValid = true
				}

				if !localValid && !hashValidated {
					// Next Attempt to flip hash bits as well

					for _, hashBits := range bits {
						hash.swapLeastConfidentBitsWithBits(hashBits)

						if data.verifyHash(hash) {
							localValid = true
							break
						}
					}
				}

				if localValid {
					d.data = data.data
					d.confidence = data.confidence

					valid = true
				}

			}
			//})
		}

		if valid {
			break
		}
	}

	return valid
}

func (data DataConfidenceCollection) verifyHash(hash *DataConfidenceCollection) bool {
	if data.HashVerifier != nil {
		return data.HashVerifier(data.data, hash.data)
	}

	equal := true

	for i, value := range data.data {
		if value != hash.data[i] {
			equal = false
			break
		}
	}

	return equal
}

func (data *DataConfidenceCollection) swapLeastConfidentBitsWithBits(bits []bool) {
	bitLength := len(bits)

	swapIndexes := data.LeastConfidentBitIndexes(bitLength)

	for bitIndex, dataIndex := range swapIndexes {
		data.data[dataIndex] = bits[bitIndex]
	}
}

func possibleBinarySlicesWithNumberOfBits(bitCount int) [][]bool {
	sliceSize := int(math.Pow(2, float64(bitCount)))
	offset := 0

	if bitCount > 1 {
		offset = int(math.Pow(2, float64(bitCount-1)))
		sliceSize -= offset
	}

	possiblities := make([][]bool, sliceSize)

	for i := range possiblities {
		bin := strconv.FormatInt(int64(i+offset), 2)
		bin = utils.LeftPad(bin, bitCount)

		boolBin := make([]bool, bitCount)

		for index, value := range bin {
			boolBin[index] = value == '1'
		}

		possiblities[i] = boolBin
	}

	return possiblities
}
