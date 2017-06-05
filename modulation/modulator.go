package modulation

import (
	"fmt"
	"math"
	"math/rand"

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

	CachedWave *[]float64
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
	if m.CachedWave == nil {
		numberOfFrames := int64(m.FrameCount())

		sampleRate := int64(m.Options.Kilobitrate * 1000)
		numberOfSamplesPerFrame := int64(float64(sampleRate) * (float64(m.Options.FrameDuration) / float64(1000.0)))

		offsetSamples := numberOfSamplesPerFrame
		numberOfSamples := (numberOfSamplesPerFrame * numberOfFrames) + offsetSamples

		fullWave := make([]float64, numberOfSamples)

		carrierData := make([][]float64, m.carrierCount)

		ts := 1 / float64(sampleRate)

		for i := int(0); i < int(numberOfFrames); i++ {
			frame := m.buildFrame(i)

			for carrierIndex, freq := range frame.allCarriers() {
				carrierData[carrierIndex] = append(carrierData[carrierIndex], freq)
			}
		}

		scaler := float64(math.Pow(2, float64(options.BitDepth-2)))
		signalScaler := scaler

		for phase := int64(0); phase < numberOfSamples; phase++ {
			amplitude := float64(0.0)
			carrierCount := 0

			for carrierIndex, frequencies := range carrierData {
				offsetIndex := int64(carrierIndex % options.OffsetGroups)
				carrierOffset := int64(0)

				if offsetIndex != 0 {
					carrierOffset = int64(offsetSamples / offsetIndex)
				}

				frequencyIndex := phase - carrierOffset

				if frequencyIndex >= 0 && frequencyIndex < (numberOfSamplesPerFrame*numberOfFrames) {
					frameIndex := frequencyIndex / numberOfSamplesPerFrame
					frequency := frequencies[frameIndex]

					if frequency != -1.0 {
						amplitude += (math.Sin(float64(phase) * ts * frequency * 2 * (math.Pi)))
						carrierCount++
					}
				}
			}

			if carrierCount != 0 {
				amplitude = amplitude / float64(carrierCount)
			}

			noise := float64(0)
			scaler := float64(signalScaler)

			if m.Options.NoiseLevel != 0 {
				noiseAmplitude := (scaler / float64(100.0)) * float64(m.Options.NoiseLevel)

				scaler = scaler - noiseAmplitude
				noise = noiseAmplitude * rand.Float64()
			}

			fullWave[phase] = (amplitude * scaler) + noise
		}

		fmt.Println("wave generated")

		m.CachedWave = &fullWave
	}

	return *m.CachedWave
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
	if frameIndex <= 1 {
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
