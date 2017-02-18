package modulation

import (
	"fmt"
	"math"

	"github.com/Hunter-Dolan/midrange/options"
)

type Modulator struct {
	Options *options.Options

	carrierCount   int
	carrierSpacing float64

	data *[]bool

	dataLength int

	frameCount int
	frameIndex int
}

func (m *Modulator) SetData(data *[]bool) {
	m.data = data
	m.setup()
}

func (m *Modulator) NextWave() *[]float64 {
	frame := m.buildFrame(m.frameIndex)
	m.frameIndex++

	return frame.wave()
}

func (m *Modulator) FullWave() []float64 {

	numberOfFrames := int64(m.FrameCount())

	sampleRate := int64(m.Options.Kilobitrate * 1000)
	numberOfSamplesPerFrame := int64(float64(sampleRate) * (float64(m.Options.FrameDuration) / float64(1000.0)))
	numberOfSamples := numberOfSamplesPerFrame * numberOfFrames

	fullWave := make([]float64, numberOfSamples)

	for waveIndex := int64(0); waveIndex < numberOfFrames; waveIndex++ {
		wave := *m.NextWave()

		for i, amplitude := range wave {
			fullWave[waveIndex*numberOfSamplesPerFrame+int64(i)] = amplitude
		}
	}

	return fullWave
}

func (m *Modulator) Reset() {
	m.frameIndex = 0
}

func (m *Modulator) FrameCount() int {
	return m.frameCount
}

func (m *Modulator) setup() {
	m.carrierSpacing = float64(m.Options.OMFSKConstant) / (float64(m.Options.FrameDuration) / 1000.0)
	m.carrierCount = int(math.Floor(float64(m.Options.Bandwidth) / m.carrierSpacing))
	m.dataLength = len(*m.data)
	m.frameCount = int(math.Ceil(float64(m.dataLength)/float64(m.carrierCount)) + 2)
	m.frameIndex = 0

	fmt.Println((m.carrierCount * (1000 / m.Options.FrameDuration)), "baud")
}

func (m *Modulator) buildHeaderFrame(frameIndex int) frame {
	data := make([]bool, m.carrierCount)

	for i := 0; i < m.carrierCount; i++ {
		data[i] = i%2 == frameIndex
	}

	f := frame{}
	f.data = data
	f.carrierSpacing = m.carrierSpacing
	f.frameIndex = frameIndex
	f.options = m.Options
	f.header = true

	return f
}

func (m Modulator) buildFrame(frameIndex int) frame {
	if m.frameIndex <= 1 {
		return m.buildHeaderFrame(frameIndex)
	}

	leftEnd := m.carrierCount * (frameIndex - 2)
	rightEnd := int(leftEnd + m.carrierCount)

	if rightEnd > m.dataLength {
		rightEnd = m.dataLength - 1
	}

	data := (*m.data)[leftEnd:rightEnd]

	f := frame{}
	f.data = data
	f.carrierSpacing = m.carrierSpacing
	f.frameIndex = frameIndex
	f.options = m.Options

	return f

}
