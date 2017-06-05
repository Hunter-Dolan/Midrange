package transaction

import (
	"fmt"
	"math"

	"github.com/Hunter-Dolan/midrange/frame"
)

// Transaction is the base transaction structure
type Transaction struct {
	Frames        []*frame.Frame
	BaseFrequency int
	FrameDuration int
	Kilobitrate   int
	Bandwidth     int

	OMFSKConstant float64

	carrierCount   int
	carrierSpacing float64

	wave []float64

	// Debug
	NoiseLevel int
}

// NewTransaction creates a new transaction
func NewTransaction() *Transaction {
	t := Transaction{}
	t.BaseFrequency = 50000
	t.FrameDuration = 500
	t.Kilobitrate = 96 * 2
	t.Bandwidth = 1000
	t.OMFSKConstant = 1.0
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

func (t *Transaction) determineCarrierCount() {
	t.carrierSpacing = float64(t.OMFSKConstant) / (float64(t.FrameDuration) / 1000.0)
	t.carrierCount = int(math.Floor(float64(t.Bandwidth) / t.carrierSpacing))
}

// SetData sets the data for the transaction
func (t *Transaction) SetData(s string) {

	t.determineCarrierCount()

	bin := stringToBin(s)
	binLength := len(bin)

	frameData := []int{}

	frameSum := 0

	for i, binaryBit := range bin {
		carrierIndex := i % t.carrierCount
		bit := int(binaryBit) - 48
		frameData = append(frameData, bit)
		frameSum += bit

		if carrierIndex == (t.carrierCount-1) || i == (binLength-1) {
			byteLength := len(frameData) / 8
			byteOffset := len(t.Frames) * byteLength

			fmt.Println(frameData, s[byteOffset:byteLength+byteOffset], frameSum)

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
	frameOptions := frame.GenerationOptions{}
	frameOptions.Duration = int64(t.FrameDuration)
	frameOptions.BaseFrequency = t.BaseFrequency
	frameOptions.SampleRate = int64(t.Kilobitrate * 1000)
	frameOptions.CarrierCount = t.carrierCount
	frameOptions.CarrierSpacing = t.carrierSpacing
	frameOptions.NoiseLevel = t.NoiseLevel

	return &frameOptions
}

func (t *Transaction) Wave() []float64 {
	if t.wave == nil {
		t.buildHeader()

		wave := []float64{}

		fmt.Println(float64(t.carrierCount)*(1000.0/float64(t.FrameDuration)), "Bits/second")
		fmt.Println(t.carrierCount, "Carriers")

		numFrames := len(t.Frames)

		waveIndex := int64(0)

		frameOptions := t.FrameGenerationOptions()

		for frameIndex := 0; frameIndex < numFrames; frameIndex++ {
			frame := t.Frames[frameIndex]
			waveIndex = frame.Generate(frameOptions, waveIndex)
			wave = append(wave, frame.Wave...)
		}

		fmt.Println(numFrames, "Frames")

		t.wave = wave
	}

	return t.wave
}
