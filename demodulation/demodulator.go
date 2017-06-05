package demodulation

import (
	"crypto/sha1"
	"fmt"
	"math"
	"os"

	"github.com/Hunter-Dolan/midrange/fec"
	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

type Demodulator struct {
	Options *options.Options

	carrierCount   int
	carrierSpacing float64

	frames               []*frame
	confidenceCollection fec.DataConfidenceCollection

	wave       *[]float64
	waveLength int

	trainer trainer

	dataLength int

	header  *header
	packets []*packet
}

func (d *Demodulator) SetWave(wave *[]float64) {
	d.wave = wave
	d.setup()
}

func (d *Demodulator) Data() []byte {
	data := []byte{}

	for _, packet := range d.packets {
		dataBytes := utils.BinDecode(packet.data.Data())
		data = append(data, dataBytes...)
	}

	data = decompress(data)

	return data
}

func (d *Demodulator) setup() {
	d.carrierSpacing = float64(d.Options.OMFSKConstant) / (float64(d.Options.FrameDuration) / 1000.0)
	d.carrierCount = int(math.Floor(float64(d.Options.Bandwidth) / d.carrierSpacing))
	d.waveLength = len(*d.wave)

	d.train()
	d.demodulateFrames()

	d.buildHeader()
	d.buildPackets()
	d.correctPacketHashes()
	d.correctPackets()
}

func (d *Demodulator) train() {
	d.trainer = trainer{}
	d.trainer.demodulator = d
	d.trainer.train()
}

func (d *Demodulator) demodulateFrames() {
	fmt.Println("Demodulating Frames")

	sampleRate := int64(d.Options.Kilobitrate * 1000)
	numberOfSamplesPerFrame := int(float64(sampleRate) * (float64(d.Options.FrameDuration) / float64(1000)))

	numberOfFrames := int(math.Floor(float64(d.waveLength/numberOfSamplesPerFrame))) - 2

	d.frames = make([]*frame, numberOfFrames)

	d.confidenceCollection = fec.DataConfidenceCollection{}

	for index := range d.frames {
		f := frame{}
		f.index = index + 2
		f.demodulator = d
		d.frames[index] = &f
		f.demodulate()
		d.confidenceCollection.Append(f.confidenceCollection)
	}

	d.dataLength = d.confidenceCollection.Length()

	d.confidenceCollection.Deinterleave()
}

func (d *Demodulator) buildHeader() {
	fmt.Println("Decoding Header")

	confidenceCollection := d.confidenceCollection.Slice(0, options.HeaderLengthWithoutRedundancy)
	hashCollection := d.confidenceCollection.Slice(options.HeaderLengthBits, options.PacketChecksumLength)

	for i := 0; i < options.HeaderChecksumRedundancy; i++ {
		offset := options.HeaderLengthWithoutRedundancy + (i * options.PacketChecksumLength)
		redundantChecksum := d.confidenceCollection.Slice(offset, options.PacketChecksumLength)

		hashCollection.AddRedundantCollection(redundantChecksum)
	}

	hashCollection.CorrectDataWithRedundancy()

	d.header = &header{}
	d.header.demodulator = d
	d.header.data = confidenceCollection
	d.header.hash = hashCollection

	d.header.correct()
	d.header.parse()
}

func (d *Demodulator) buildPackets() {
	fmt.Println("Decoding Packets")

	headerLength := d.header.Length()
	packetCount := (d.confidenceCollection.Length() - headerLength) / options.PacketTotalLength

	fmt.Println(packetCount, (d.confidenceCollection.Length() - headerLength))

	d.packets = make([]*packet, packetCount)

	for i := range d.packets {
		offset := headerLength + (i * options.PacketTotalLength)
		packetTotalData := d.confidenceCollection.Slice(offset, options.PacketTotalLength)

		packet := packet{}

		dataLength := options.PacketDataLength

		if i == packetCount-1 {
			dataLength = d.header.dataLength - (i * options.PacketDataLength)
		}

		fmt.Println(d.header.dataLength, i*options.PacketDataLength)

		packet.data = packetTotalData.Slice(0, dataLength)

		packet.hash = packetTotalData.Slice(dataLength, options.PacketChecksumLength)

		d.packets[i] = &packet
	}

}

func (d *Demodulator) correctPacketHashes() {
	fmt.Println("Decoding Packet Hashes")

	hashCollection := fec.DataConfidenceCollection{}

	for _, packet := range d.packets {
		hashCollection.Append(packet.hash)
	}

	hash := d.header.transmissionHash

	hashCollection.HashVerifier = func(data []bool, hash []bool) bool {

		dataByte := utils.BoolToByteBin(data)

		sha1DataFull := sha1.Sum(dataByte)
		sha1Data := sha1DataFull[:options.HeaderTransmissionChecksumLength/8]

		boolBinHash := utils.BinEncodeWithPad(sha1Data, options.HeaderTransmissionChecksumLength)

		match := true

		for i, value := range hash {
			if boolBinHash[i] != value {
				match = false
				break
			}
		}

		return match
	}

	if !hashCollection.CorrectData(hash, true) {
		fmt.Println("Could not decode packet hashes")
		os.Exit(0)
	}

	for i, packet := range d.packets {
		packet.hash = hashCollection.Slice(i*options.PacketChecksumLength, options.PacketChecksumLength)
	}

}

func (d *Demodulator) correctPackets() {
	fmt.Println("Correcting Packets")

	for _, packet := range d.packets {
		packet.Correct()
	}
}
