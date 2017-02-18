package demodulation

import (
	"crypto/sha1"

	"github.com/Hunter-Dolan/midrange/fec"
	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

type packet struct {
	data *fec.DataConfidenceCollection
	hash *fec.DataConfidenceCollection
}

func (p *packet) Correct() bool {

	p.data.HashVerifier = func(data []bool, hash []bool) bool {
		dataByte := utils.BinDecode(data)

		sha1DataFull := sha1.Sum(dataByte)
		sha1Data := sha1DataFull[:options.PacketChecksumLength/8]

		boolBinHash := utils.BinEncodeWithPad(sha1Data, options.PacketChecksumLength)

		match := true

		for i, value := range hash {
			if boolBinHash[i] != value {
				match = false
				break
			}
		}

		return match
	}

	return p.data.CorrectData(p.hash, true)
}

func (p *packet) Length() int {
	return p.data.Length() + p.hash.Length()
}
