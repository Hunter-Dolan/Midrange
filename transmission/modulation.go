package transmission

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"github.com/Hunter-Dolan/midrange/modulation"
	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

func (t *Transmission) SetData(data string) {
	dataBytes := []byte(data)

	dataBytes = compress(dataBytes)

	t.dataLength = len(dataBytes) * 8
	t.originalDataLength = len(data) * 8

	packetCount := math.Ceil(float64(t.dataLength) / float64(options.PacketDataLength))

	t.data = make([]*packet, int(packetCount))

	for packetIndex := range t.data {
		p := packet{}

		leftEnd := int(packetIndex * options.PacketDataLengthBytes)
		rightEnd := int(leftEnd + options.PacketDataLengthBytes)

		if rightEnd > len(dataBytes) {
			rightEnd = len(dataBytes)
		}

		p.setBytes(dataBytes[leftEnd:rightEnd])

		t.data[packetIndex] = &p
	}

	t.generateHeaderPacket()
}

func (t *Transmission) generateHeaderPacket() {
	dataLengthBin, dataLengthBool := t.dataLengthBin()
	checksumBin, checksumBool := t.checksumBin()

	bytes := []byte((dataLengthBin + checksumBin)[:])

	bin := append(dataLengthBool, checksumBool...)

	p := packet{}
	p.data = &bin

	p.calculateHash(bytes)
	p.replicateHash(options.HeaderChecksumRedundancy)

	t.header = &p
}

func (t *Transmission) checksumBin() (string, []bool) {
	checksum := []byte{}

	for _, packet := range t.data {
		checksum = append(checksum, packet.hashBytes...)
	}

	checksumBin := utils.BinEncode(checksum)

	checksumBinByte := utils.BoolToByteBin(checksumBin)

	sha1DataFull := sha1.Sum(checksumBinByte)
	sha1Data := sha1DataFull[:options.HeaderTransmissionChecksumLength/8]

	sha1Int, _ := strconv.ParseInt(hex.EncodeToString(sha1Data), 16, 64)

	sha1Bin := strconv.FormatInt(sha1Int, 2)
	sha1Bin = utils.LeftPad(sha1Bin, options.HeaderTransmissionChecksumLength)

	sha1Dec := strconv.FormatInt(sha1Int, 10)
	fmt.Println("Transmission Hash:", sha1Dec)

	return sha1Bin, utils.BoolBinEncode(sha1Bin)
}

func (t *Transmission) dataLengthBin() (string, []bool) {
	dataLength := 0

	for _, packet := range t.data {
		dataLength += len(*packet.data)
	}

	t.dataLength = dataLength

	headerDataLengthInteger := int64(dataLength)
	headerDataLengthBin := strconv.FormatInt(headerDataLengthInteger, 2)
	headerDataLengthBin = utils.LeftPad(headerDataLengthBin, options.HeaderDataLengthIntegerBitLength)

	headerDataLengthDec := strconv.FormatInt(headerDataLengthInteger, 10)
	fmt.Println("Data Length:", headerDataLengthDec)

	return headerDataLengthBin, utils.BoolBinEncode(headerDataLengthBin)
}

func (t *Transmission) allData() *[]bool {
	data := t.header.packetData()

	for _, packet := range t.data {
		data = append(data, packet.packetData()...)
	}

	carrierSpacing := float64(t.options.OMFSKConstant) / (float64(t.options.FrameDuration) / 1000.0)
	carrierCount := int(math.Floor(float64(t.options.Bandwidth) / carrierSpacing))

	fullLength := int(math.Ceil(float64(len(data))/float64(carrierCount))) * carrierCount

	paddingAmount := fullLength - len(data)

	for i := 0; i < paddingAmount; i++ {
		data = append(data, false)
	}

	data = interleave(data)

	return &data
}

func (t *Transmission) modulate() {
	if t.modulator == nil {
		t.modulator = &modulation.Modulator{}
		t.modulator.Options = t.options

		allData := t.allData()
		t.modulator.SetData(allData)
	}

	t.modulator.Reset()

	duration := float64(t.modulator.FrameCount()*t.options.FrameDuration) / 1000

	fmt.Println(t.originalDataLength, "bits to transfer")
	fmt.Println(duration, "second transfer time")

	fmt.Println(int64(float64(t.dataLength)/duration), "bps real")
	fmt.Println(int64(float64(t.originalDataLength)/duration), "bps effective")
}

func (t *Transmission) BuildWav(filename string) {
	t.modulate()
	t.modulator.BuildWav(filename)
}

func (t *Transmission) Wave() []float64 {
	t.modulate()
	return t.modulator.FullWave()
}
