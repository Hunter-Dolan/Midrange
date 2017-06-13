package demodulation

import (
	"math"

	"github.com/Hunter-Dolan/midrange/options"
	"github.com/mjibson/go-dsp/spectral"
)

/*
// A test function to determine the best nftt value
func (d *Demodulator) carrierPowerAtFrame(frameIndex int) []float64 {
	bestValue := 9999999999999.0
	offset := 0.0001
	start := 16.837
	iteration := 0
	best := 0.000

	// 53.70708024275132

	for {
		testValue := start + (offset * float64(iteration))

		fmt.Println("Testing:", testValue)

		powers := make([][]float64, options.OffsetGroups)
		frequencies := make([][]float64, options.OffsetGroups)

		sampleRate := int64(d.Options.Kilobitrate * 1000)
		numberOfSamplesPerFrame := int(float64(sampleRate) * (float64(d.Options.FrameDuration) / float64(1000.0)))

		for offsetIndex := range powers {
			offsetAmount := 0

			if options.OffsetGroups != 1 {
				offsetAmount = (numberOfSamplesPerFrame / int(options.OffsetGroups)) * offsetIndex
			}

			leftEnd := int(frameIndex*numberOfSamplesPerFrame) + offsetAmount

			guard := numberOfSamplesPerFrame - (numberOfSamplesPerFrame / options.GuardDivisor)

			rightEnd := (int(leftEnd+numberOfSamplesPerFrame) + offsetAmount) - guard

			waveLength := d.waveLength

			//fmt.Println(frameIndex, offsetAmount, leftEnd, rightEnd, waveLength)

			if rightEnd > waveLength {
				rightEnd = waveLength - 1
			}

			segment := (*d.wave)[leftEnd:rightEnd]

			//nfft := math.Pow(2, d.Options.NFFTPower)
			//	nfft := math.Pow(2, testValue)

			// Perform Pwelch on wave
			fs := float64(d.Options.Kilobitrate * 1000)
			opts := &spectral.PwelchOptions{
				//NFFT: int(math.Pow(2, 17.0)),
				NFFT: int(math.Pow(2, testValue)),
			}

			offsetPowers, offsetFrequencies := spectral.Pwelch(segment, fs, opts)
			powers[offsetIndex] = offsetPowers
			frequencies[offsetIndex] = offsetFrequencies
		}

		carriers := d.carrierFrequencies()
		numberOfCarriers := len(carriers)

		minimumCarrierAdjusted := float64(carriers[0] - 5)
		maximumCarrierAdjusted := float64(carriers[numberOfCarriers-1] + 5)

		carriersFound := 0

		locatedCarrierValues := make([]float64, numberOfCarriers)

		totalDist := 0.0

		for i, searchingCarrier := range carriers {
			groupIndex := i % options.OffsetGroups
			groupFrequencies := frequencies[groupIndex]

			minimumDistance := float64(-1)

			for i, frequency := range groupFrequencies {
				if frequency < maximumCarrierAdjusted && frequency > minimumCarrierAdjusted {

					distance := math.Abs(frequency - float64(searchingCarrier))

					if minimumDistance < distance && minimumDistance != float64(-1) {
						locatedCarrierValues[carriersFound] = powers[groupIndex][i-1]

						totalDist += math.Abs(searchingCarrier - groupFrequencies[i-1])

						carriersFound++
						break
					} else {
						minimumDistance = distance
					}
				}
			}
		}

		fmt.Println(totalDist)

		if bestValue > totalDist {
			fmt.Println("So far Best:", start+(offset*float64(iteration-1)))
			//break
			best = start + (offset * float64(iteration-1))
			bestValue = totalDist
		}

		iteration++

		if iteration > 100 {
			break
		}
	}

	fmt.Println("Aboslute Best:", best, bestValue)

	os.Exit(0)

	return []float64{}
}

/*/

func (d *Demodulator) carrierPowerAtFrame(frameIndex int) []float64 {

	powers := make([][]float64, options.OffsetGroups)
	frequencies := make([][]float64, options.OffsetGroups)

	sampleRate := int64(d.Options.Kilobitrate * 1000)
	numberOfSamplesPerFrame := int(float64(sampleRate) * (float64(d.Options.FrameDuration) / float64(1000.0)))

	for offsetIndex := range powers {
		offsetAmount := 0

		if options.OffsetGroups != 1 {
			offsetAmount = (numberOfSamplesPerFrame / int(options.OffsetGroups)) * offsetIndex
		}

		leftEnd := int(frameIndex*numberOfSamplesPerFrame) + offsetAmount

		guard := numberOfSamplesPerFrame - (numberOfSamplesPerFrame / options.GuardDivisor)

		rightEnd := (int(leftEnd+numberOfSamplesPerFrame) + offsetAmount) - guard

		waveLength := d.waveLength

		//fmt.Println(frameIndex, offsetAmount, leftEnd, rightEnd, waveLength)

		if rightEnd > waveLength {
			rightEnd = waveLength - 1
		}

		segment := (*d.wave)[leftEnd:rightEnd]

		nfft := math.Pow(2, d.Options.NFFTPower)
		pad := math.Pow(2, d.Options.NFFTPower)

		// Perform Pwelch on wave
		fs := float64(d.Options.Kilobitrate * 1000)
		opts := &spectral.PwelchOptions{
			NFFT: int(nfft),
			Pad:  int(pad),
		}

		offsetPowers, offsetFrequencies := spectral.Pwelch(segment, fs, opts)
		powers[offsetIndex] = offsetPowers
		frequencies[offsetIndex] = offsetFrequencies
	}

	carriers := d.carrierFrequencies()
	numberOfCarriers := len(carriers)

	minimumCarrierAdjusted := float64(carriers[0] - 5)
	maximumCarrierAdjusted := float64(carriers[numberOfCarriers-1] + 5)

	carriersFound := 0

	locatedCarrierValues := make([]float64, numberOfCarriers)

	for i, searchingCarrier := range carriers {
		groupIndex := i % options.OffsetGroups
		groupFrequencies := frequencies[groupIndex]

		minimumDistance := float64(-1)

		for i, frequency := range groupFrequencies {
			if frequency < maximumCarrierAdjusted && frequency > minimumCarrierAdjusted {

				distance := math.Abs(frequency - float64(searchingCarrier))

				if minimumDistance < distance && minimumDistance != float64(-1) {
					locatedCarrierValues[carriersFound] = powers[groupIndex][i-1]
					//fmt.Println(groupFrequencies[i-1], searchingCarrier)
					carriersFound++
					break
				} else {
					minimumDistance = distance
				}
			}
		}
	}

	//fmt.Println(locatedCarrierValues)

	return locatedCarrierValues
}

//*/

func (d *Demodulator) carrierFrequencies() []float64 {
	var carriers = make([]float64, d.carrierCount)

	for index := range carriers {
		freq := float64(index)*d.carrierSpacing + float64(d.Options.BaseFrequency)
		carriers[index] = freq
	}

	return carriers
}
