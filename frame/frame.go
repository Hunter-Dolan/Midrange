package frame

import "math"

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
}

// NewFrame creates a new frame
func NewFrame(data ...int) *Frame {
	frame := Frame{}
	frame.Data = data
	frame.SignalFrequency = -1

	return &frame
}

// NewHeaderFrame creates a header frame
func NewHeaderFrame(options *GenerationOptions) *Frame {
	data := make([]int, options.CarrierCount)
	data[0] = 1
	data[options.CarrierCount-1] = 1

	header := NewFrame(data...)
	header.SignalFrequency = 300

	return header
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

	for i := int64(0); i < numSamples; i++ {

		amplitude := float64(0)

		p := float64(float64(i+startIndex) * ts)

		for carrierIndex := int64(0); carrierIndex < int64(carrierCount); carrierIndex++ {
			freq := float64(carriers[carrierIndex])
			amplitude += (math.Sin(p * freq * 2 * (math.Pi)))
		}

		amplitude = (amplitude / float64(carrierCount))

		f.Wave = append(f.Wave, amplitude*32767.0) //+(rand.Float64()*5000.0))

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
