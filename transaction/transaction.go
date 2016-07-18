package transaction

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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
}

// NewTransaction creates a new transaction
func NewTransaction() *Transaction {
	t := Transaction{}
	t.BaseFrequency = 10000
	t.FrameDuration = 500
	t.Carriers = 64
	t.Kilobitrate = 96 * 2
	t.Bandwidth = 1000

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
	/*
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		w.Write([]byte(s))
		w.Close()

		s = b.String()
	*/
	bin := stringToBin(s)
	binLength := len(bin)

	frameData := []int{}

	for i, bit := range bin {
		frameIndex := i % t.Carriers
		frameData = append(frameData, int(bit)-48)

		if frameIndex == (t.Carriers-1) || i == (binLength-1) {
			f := frame.NewFrame(frameData...)
			t.AddFrame(f)
			frameData = []int{}
		}
	}

}

func (t *Transaction) buildHeader() {
	header := frame.NewHeaderFrame(t.FrameGenerationOptions())

	t.Frames = append([]*frame.Frame{header}, t.Frames...)
}

func (t Transaction) FrameGenerationOptions() *frame.GenerationOptions {
	spacing := t.Bandwidth / t.Carriers

	frameOptions := frame.GenerationOptions{}
	frameOptions.Duration = int64(t.FrameDuration)
	frameOptions.BaseFrequency = t.BaseFrequency
	frameOptions.SampleRate = int64(t.Kilobitrate * 1000)
	frameOptions.Spacing = spacing
	frameOptions.CarrierCount = t.Carriers

	return &frameOptions
}

func (t *Transaction) Wave() []float64 {
	t.buildHeader()

	wave := []float64{}

	fmt.Println(float64(t.Carriers)*(1000.0/float64(t.FrameDuration)), "Bits/second")

	numFrames := len(t.Frames)

	waveIndex := int64(0)

	frameOptions := t.FrameGenerationOptions()

	for frameIndex := 0; frameIndex < numFrames; frameIndex++ {
		frame := t.Frames[frameIndex]
		waveIndex = frame.Generate(frameOptions, waveIndex)
		wave = append(wave, frame.Wave...)
	}

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
