package demodulation

import (
	"math"

	"github.com/mjibson/go-dsp/spectral"
)

func (d *Demodulator) carrierPowerAtFrame(frameIndex int) []float64 {

	sampleRate := int64(d.Options.Kilobitrate * 1000)
	numberOfSamplesPerFrame := int(float64(sampleRate) * (float64(d.Options.FrameDuration) / float64(1000.0)))

	leftEnd := int(frameIndex * numberOfSamplesPerFrame)
	rightEnd := int(leftEnd + numberOfSamplesPerFrame)
	waveLength := d.waveLength

	if rightEnd > waveLength {
		rightEnd = waveLength - 1
	}

	segment := (*d.wave)[leftEnd:rightEnd]

	nfft := math.Pow(2, float64(d.Options.NFFTPower))

	// Perform Pwelch on wave
	fs := float64(d.Options.Kilobitrate * 1000)
	opts := &spectral.PwelchOptions{
		NFFT: int(nfft),
	}

	powers, frequencies := spectral.Pwelch(segment, fs, opts)

	carriers := d.carrierFrequencies()
	numberOfCarriers := len(carriers)

	minimumCarrierAdjusted := float64(carriers[0] - 5)
	maximumCarrierAdjusted := float64(carriers[numberOfCarriers-1] + 5)

	carriersFound := 0
	minimumDistance := float64(-1)
	searchingCarrier := carriers[0]

	locatedCarrierValues := make([]float64, numberOfCarriers)

	for i, frequency := range frequencies {
		if frequency < maximumCarrierAdjusted && frequency > minimumCarrierAdjusted {

			distance := math.Abs(frequency - float64(searchingCarrier))

			if minimumDistance < distance && minimumDistance != float64(-1) {
				locatedCarrierValues[carriersFound] = powers[i-1]
				carriersFound++

				if carriersFound < numberOfCarriers {
					searchingCarrier = carriers[carriersFound]
					minimumDistance = float64(-1)
				} else {
					break
				}
			} else {
				minimumDistance = distance
			}
		}
	}

	return locatedCarrierValues
}

func (d *Demodulator) carrierFrequencies() []float64 {
	var carriers = make([]float64, d.carrierCount)

	for index := range carriers {
		freq := float64(index)*d.carrierSpacing + float64(d.Options.BaseFrequency)
		carriers[index] = freq
	}

	return carriers
}
