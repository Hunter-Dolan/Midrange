package frame

import (
	"math"
	"math/rand"
	"strconv"
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
	KeyStates     int

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
	evenFrame.HeaderPacket = true

	oddFrame := NewFrame(oddData...)
	oddFrame.HeaderPacket = true

	return []*Frame{evenFrame, oddFrame}
}

func (f Frame) carriers(options *GenerationOptions) map[float64]int {
	var carriers = map[float64]int{}

	var dataLen = len(f.Data)

	bitsPerCarrier := int(math.Log2(float64(options.KeyStates)))

	activeCarrierIndex := 0
	inactiveCarrierIndex := 1

	activeCarrier := float64(activeCarrierIndex*options.Spacing + options.BaseFrequency)
	carriers[activeCarrier] = options.KeyStates - 1

	inactiveCarrier := float64(inactiveCarrierIndex*options.Spacing + options.BaseFrequency)
	carriers[inactiveCarrier] = 0

	if f.HeaderPacket {
		signalFrequency := 300.0
		carriers[signalFrequency] = 1

		for index, bit := range f.Data {
			freq := float64(index*options.Spacing + options.BaseFrequency)

			if bit == 1 {
				carriers[freq] = options.KeyStates - 1
			}
		}

	} else {
		for index := 0; index < dataLen/bitsPerCarrier; index++ {

			offset := index * bitsPerCarrier
			rightEnd := offset + bitsPerCarrier

			if rightEnd > dataLen {
				rightEnd = dataLen - 1
			}

			segment := f.Data[offset:rightEnd]

			stringValue := ""

			for _, s := range segment {
				stringValue += strconv.Itoa(s)
			}

			value, _ := strconv.ParseInt(stringValue, 2, 64)
			intValue := int(value)

			if intValue != 0 {
				freq := float64((index+2)*options.Spacing + options.BaseFrequency)
				carriers[freq] = intValue
			}
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

		for freq, state := range carriers {
			freqAmp := float64(float64(state) / float64(options.KeyStates-1))

			amplitude += (math.Sin(p * freq * 2 * math.Pi)) * freqAmp
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
