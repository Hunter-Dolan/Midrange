package demodulation

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Hunter-Dolan/midrange/fec"
	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

type header struct {
	packet

	demodulator *Demodulator

	transmissionHash *fec.DataConfidenceCollection

	dataLength int
}

func (h *header) correct() bool {
	h.data.HashVerifier = func(data []bool, hash []bool) bool {
		valid := fec.SHA1DirectBinVerifier(data, hash)

		if valid {
			intLength := options.HeaderDataLengthIntegerBitLength
			dataLengthBool := data[:intLength]
			dataLengthBin := utils.BoolToBin(dataLengthBool)
			dataLengthDec, _ := strconv.ParseInt(dataLengthBin, 2, 64)
			dataLength := int(dataLengthDec)

			approxBitsMax := ((h.demodulator.dataLength - options.HeaderLengthBits) / options.PacketTotalLength) * options.PacketDataLength
			approxBitsMin := approxBitsMax - h.demodulator.carrierCount

			if dataLength <= approxBitsMax && dataLength > approxBitsMin {
				return true
			}
		}

		return false
	}

	correct := h.data.CorrectData(h.hash, true)

	if !correct {
		correct = h.data.CorrectData(h.hash, false)

		if !correct {
			fmt.Println("Unable to decode header")
			os.Exit(0)
		}
	}

	return correct
}

func (h *header) parse() {
	boolData := h.data.Data()

	intLength := options.HeaderDataLengthIntegerBitLength
	checkLength := options.HeaderTransmissionChecksumLength

	dataLengthBool := boolData[:intLength]
	dataLengthBin := utils.BoolToBin(dataLengthBool)
	dataLengthDec, _ := strconv.ParseInt(dataLengthBin, 2, 64)
	h.dataLength = int(dataLengthDec)

	h.transmissionHash = h.data.Slice(intLength, checkLength)

}

func (p *header) Length() int {
	return options.HeaderLengthBits + options.PacketChecksumLength
}
