package fec

import (
	"crypto/sha1"

	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

func SHA1DirectBinVerifier(data []bool, hash []bool) bool {
	dataByte := utils.BoolToByteBin(data)

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
