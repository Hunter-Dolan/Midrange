package matcher

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Hunter-Dolan/midrange/frame"
	"github.com/mjibson/go-dsp/spectral"
)

// FindProbableMatch returns the most likely match for a wave
func (b *Matcher) findProbableMatch(wave []float64) *frame.Frame {

	frame := &frame.Frame{}

	nfft := math.Log(float64(b.options.SampleRate * 1000))
	nfft = math.Ceil(nfft / math.Log(2))

	nfft = math.Pow(2, float64(b.options.NFFTPower))

	//topCandidate := b.frames[0]

	// Perform Pwelch on wave
	fs := float64(b.options.SampleRate)
	opts := &spectral.PwelchOptions{
		NFFT:      int(nfft),
		Scale_off: false,
	}

	powers, frequencies := spectral.Pwelch(wave, fs, opts)

	carriers := b.options.Carriers()
	numberOfCarriers := len(carriers)

	minimumCarrierAdjusted := float64(carriers[0] - 5)
	maximumCarrierAdjusted := float64(carriers[numberOfCarriers-1] + 5)

	carriersFound := 0
	minimumDistance := float64(-1)
	searchingCarrier := carriers[0]

	locatedCarrierValues := make([]float64, numberOfCarriers)

	localSamples := 2

	for i, frequency := range frequencies {

		if frequency < maximumCarrierAdjusted && frequency > minimumCarrierAdjusted {
			distance := math.Abs(frequency - float64(searchingCarrier))

			if minimumDistance < distance && minimumDistance != float64(-1) {
				step := 1
				powerLeftEnd := i - 1 - ((localSamples / 2) * step)

				midPower := powers[i-1]
				midHighest := true
				powerSum := float64(0)

				for localIndex := 0; localIndex < localSamples+1; localIndex++ {
					powerIndex := powerLeftEnd + (step * localIndex)
					power := powers[powerIndex]

					if powerIndex != (i - 1) {
						powerSum += power
					}

					if power > midPower {
						midHighest = false
					}
				}

				powerAvg := powerSum / float64(localSamples)

				aboveAverage := false

				if (midPower/powerAvg)*100 > 50 {
					aboveAverage = true
				}

				value := 0

				if midHighest && aboveAverage {
					value = 1
				}

				frame.Data = append(frame.Data, value)

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

	fmt.Println(frame.Data)

	return frame
}

// Match decodes matched frames
func (b *Matcher) match(wave []float64) []*frame.Frame {
	frames := []*frame.Frame{}

	waveLength := int64(len(wave))
	options := b.options
	frameLength := int64(float64(options.SampleRate) / float64(1000.0) * float64(options.Duration))

	numberOfFrames := waveLength / frameLength

	for i := int64(0); i < numberOfFrames; i++ {
		offset := i * frameLength
		frameWave := wave[offset : offset+frameLength]
		frame := b.findProbableMatch(frameWave)

		if len(frames) == 0 {
			frame.HeaderPacket = true
		}

		frames = append(frames, frame)
	}

	return frames
}

func (b *Matcher) Decode(wave []float64) string {
	frames := b.match(wave)
	binary := ""
	for _, frame := range frames {
		if frame.HeaderPacket == false {
			frameBinary := ""

			frameDataLen := len(frame.Data)

			padAmount := b.options.CarrierCount - frameDataLen

			for i := 0; i < padAmount; i++ {
				frameBinary += "0"
			}

			for _, bit := range frame.Data {
				frameBinary += strconv.Itoa(bit)
			}

			binary += frameBinary
		}
	}

	//fmt.Println(binary)

	bitLength := 8

	byteCount := len(binary) / bitLength
	byteArray := make([]byte, byteCount)

	for byteIndex := 0; byteIndex < byteCount; byteIndex++ {
		offset := byteIndex * bitLength
		bits := binary[offset : offset+bitLength]
		i, _ := strconv.ParseInt(bits, 2, 64)
		byteArray[byteIndex] = byte(i)
	}

	/*
		var buffer bytes.Buffer
		r, err := gzip.NewReader(&buffer)

		if err != nil {
			fmt.Println(err)
		}

		r.Read(byteArray)
		r.Close()

		return buffer.String()*/

	return string(byteArray[:])

}
