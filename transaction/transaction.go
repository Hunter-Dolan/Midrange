package transaction

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/Hunter-Dolan/midrange/frame"
)

// Transaction is the base transaction structure
type Transaction struct {
	Frames        []*frame.Frame
	BaseFrequency int
	FrameDuration int
	Carriers      int
	Kilobitrate   int
	Bandwidth     int
	KeyStates     int

	// Debug
	NoiseLevel int
}

// NewTransaction creates a new transaction
func NewTransaction() *Transaction {
	t := Transaction{}
	t.BaseFrequency = 50000
	t.FrameDuration = 500
	t.Carriers = 128
	t.Kilobitrate = 96 * 2
	t.Bandwidth = 1000
	t.KeyStates = 2

	return &t
}

// AddFrame appends a frame to the transaction
func (t *Transaction) AddFrame(f *frame.Frame) {
	t.Frames = append(t.Frames, f)
}

func stringToBin(s string) string {
	binaryString := ""

	for _, c := range s {
		binary := fmt.Sprintf("%b", c)
		binaryLen := len(binary)

		padAmount := 8 - binaryLen

		for i := 0; i < padAmount; i++ {
			binary = "0" + binary
		}

		binaryString += binary
	}
	return binaryString
}

// SetData sets the data for the transaction
func (t *Transaction) SetData(s string) {
	bin := stringToBin(s)
	binLength := len(bin)

	frameLength := (t.Carriers - 2) * (t.KeyStates / 2)

	frameData := []int{}

	frameSum := 0

	for i, binaryBit := range bin {
		bit := int(binaryBit) - 48
		frameData = append(frameData, bit)
		dataLength := len(frameData)

		frameSum += bit

		if frameLength == dataLength || i == (binLength-1) {
			byteLength := len(frameData) / 8
			byteOffset := len(t.Frames) * byteLength

			fmt.Println(frameData, s[byteOffset:byteLength+byteOffset], frameSum, len(frameData))

			f := frame.NewFrame(frameData...)
			t.AddFrame(f)
			frameData = []int{}
			frameSum = 0
		}
	}

}

func (t *Transaction) buildHeader() {
	headers := frame.NewHeaderFrames(t.FrameGenerationOptions())
	t.Frames = append(headers, t.Frames...)
}

func (t Transaction) FrameGenerationOptions() *frame.GenerationOptions {
	spacing := t.Bandwidth / t.Carriers

	frameOptions := frame.GenerationOptions{}
	frameOptions.Duration = int64(t.FrameDuration)
	frameOptions.BaseFrequency = t.BaseFrequency
	frameOptions.SampleRate = int64(t.Kilobitrate * 1000)
	frameOptions.Spacing = spacing
	frameOptions.CarrierCount = t.Carriers
	frameOptions.NoiseLevel = t.NoiseLevel
	frameOptions.KeyStates = t.KeyStates

	return &frameOptions
}

func (t *Transaction) Wave() []float64 {
	t.buildHeader()

	wave := []float64{}

	bitsPerCarrier := float64(math.Log2(float64(t.KeyStates)))

	fmt.Println(float64(t.Carriers)*(1000.0/float64(t.FrameDuration))*bitsPerCarrier, "Bits/second")

	numFrames := len(t.Frames)

	waveIndex := int64(0)

	frameOptions := t.FrameGenerationOptions()

	for frameIndex := 0; frameIndex < numFrames; frameIndex++ {
		frame := t.Frames[frameIndex]
		waveIndex = frame.Generate(frameOptions, waveIndex)
		wave = append(wave, frame.Wave...)
	}

	fmt.Println(numFrames, "Frames")

	return wave
}

// Build creates the audio for the transaction
func (t *Transaction) Build() {
	wave := t.Wave()

	rate := t.Kilobitrate * 1000
	bitDepth := 16
	duration := len(t.Frames) * t.FrameDuration

	file, err := os.Create("tone.wav")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	w := bufio.NewWriter(file)
	dataSize := rate * (duration / 1000) * (bitDepth / 8)
	buildWavHeader(w, bitDepth, rate, dataSize)

	for _, amplitude := range wave {
		binary.Write(w, binary.LittleEndian, int16(amplitude))
	}

	w.Flush()

}

const riff = "RIFF"
const wave = "WAVE"
const _fmt = "fmt "
const data = "data"

func buildWavHeader(buf io.Writer, bitDepth, sampleRate, dataSize int) {

	buf.Write([]byte(riff))
	binary.Write(buf, binary.LittleEndian, uint32(36+dataSize))
	buf.Write([]byte(wave))
	buf.Write([]byte(_fmt))
	binary.Write(buf, binary.LittleEndian, uint32(16))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))

	//ByteRate
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate*(bitDepth/8)))
	binary.Write(buf, binary.LittleEndian, uint16((bitDepth / 8)))
	binary.Write(buf, binary.LittleEndian, uint16(bitDepth))
	buf.Write([]byte(data))
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
}
