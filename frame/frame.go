package frame

import (
	"math"
	"math/rand"
)

// Frame is the base structure for a transaction frame
type Frame struct {
	Data            []int
	SignalFrequency int
	Wave            []float64
	HeaderPacket    bool
}

// GenerationOptions is the base option
type GenerationOptions struct {
	CarrierCount  int
	Duration      int64
	SampleRate    int64
	Spacing       int
	BaseFrequency int

	NoiseLevel int
}

// NewFrame creates a new frame
func NewFrame(data ...int) *Frame {
	frame := Frame{}
	frame.Data = data
	frame.SignalFrequency = -1

	return &frame
}

// NewHeaderFrames creates the header frames
func NewHeaderFrames(options *GenerationOptions) []*Frame {

	evenData := make([]int, options.CarrierCount)
	oddData := make([]int, options.CarrierCount)

	for i := 0; i < options.CarrierCount; i++ {
		even := i%2 == 0
		if even {
			evenData[i] = 1
			oddData[i] = 0
		} else {
			evenData[i] = 0
			oddData[i] = 1
		}
	}

	evenFrame := NewFrame(evenData...)
	evenFrame.SignalFrequency = 300

	oddFrame := NewFrame(oddData...)
	oddFrame.SignalFrequency = 300

	return []*Frame{evenFrame, oddFrame}
}

func (f Frame) carriers(options *GenerationOptions) []float64 {
	var carriers = []float64{}

	if f.SignalFrequency != -1 {
		carriers = append(carriers, float64(f.SignalFrequency))
	}

	for index := range f.Data {
		if f.Data[index] == 1 {
			freq := float64(index*options.Spacing + options.BaseFrequency)

			carriers = append(carriers, freq)
		}
	}

	return carriers
}

// Generate creates the wave
func (f *Frame) Generate(options *GenerationOptions, startIndex int64) int64 {
	carriers := f.carriers(options)
	carrierCount := len(carriers)

	numSamples := int64(float64(options.SampleRate) / float64(1000.0) * float64(options.Duration))

	ts := 1 / float64(options.SampleRate)

	f.Wave = make([]float64, numSamples)

	for i := int64(0); i < numSamples; i++ {

		amplitude := float64(0)

		p := float64(float64(i+startIndex) * ts)

		for carrierIndex := int64(0); carrierIndex < int64(carrierCount); carrierIndex++ {
			freq := float64(carriers[carrierIndex])
			amplitude += (math.Sin(p * freq * 2 * (math.Pi)))
		}

		amplitude = (amplitude / float64(carrierCount))

		noise := float64(0)
		scaler := float64(32767.0)

		if options.NoiseLevel != 0 {
			noiseAmplitude := (scaler / float64(100.0)) * float64(options.NoiseLevel)

			scaler = scaler - noiseAmplitude
			noise = noiseAmplitude * rand.Float64()
		}

		f.Wave[i] = amplitude*scaler + noise

	}

	return startIndex + numSamples
}

func (o *GenerationOptions) Carriers() []int {
	frequencies := make([]int, o.CarrierCount)

	for i := 0; i < o.CarrierCount; i++ {
		frequencies[i] = (i * o.Spacing) + o.BaseFrequency
	}

	return frequencies
}
