package modulation

import (
	"math"
	"math/rand"

	"github.com/Hunter-Dolan/midrange/options"
)

type frame struct {
	data []bool

	carrierSpacing float64
	frameIndex     int

	header bool

	offset int

	options *options.Options
}

func (f frame) wave() *[]float64 {
	carriers := f.carriers()
	carrierCount := len(carriers)

	sampleRate := int64(f.options.Kilobitrate * 1000)

	numSamples := int64(float64(sampleRate) / float64(1000.0) * float64(f.options.FrameDuration))

	ts := 1 / float64(sampleRate)

	wave := make([]float64, numSamples)
	startIndex := numSamples * int64(f.frameIndex)

	for i := int64(0); i < numSamples; i++ {

		amplitude := float64(0)

		if carrierCount != 0 {
			p := float64(i+startIndex) * ts

			for carrierIndex := int64(0); carrierIndex < int64(carrierCount); carrierIndex++ {
				freq := carriers[carrierIndex]
				amplitude += (math.Sin(p * freq * 2 * (math.Pi)))
			}

			amplitude = (amplitude / float64(carrierCount))
		}

		noise := float64(0)
		scaler := float64(math.Pow(2, float64(options.BitDepth-1)))

		if f.options.NoiseLevel != 0 {
			noiseAmplitude := (scaler / float64(100.0)) * float64(f.options.NoiseLevel)

			scaler = scaler - noiseAmplitude
			noise = noiseAmplitude * rand.Float64()
		}

		wave[i] = amplitude*scaler + noise

	}

	return &wave
}

func (f frame) allCarriers() []float64 {
	var carriers = []float64{}

	for index, bit := range f.data {

		if bit {
			freq := float64(index)*f.carrierSpacing + float64(f.options.BaseFrequency)
			carriers = append(carriers, freq)
		} else {
			carriers = append(carriers, -1.0)
		}
	}

	return carriers
}

func (f frame) carriers() []float64 {
	var carriers = []float64{}

	for index := range f.data {
		if f.data[index] {
			freq := float64(index)*f.carrierSpacing + float64(f.options.BaseFrequency)
			carriers = append(carriers, freq)
		}
	}

	return carriers
}
