package modulation

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Hunter-Dolan/midrange/options"
	"github.com/go-audio/aiff"
	"github.com/go-audio/audio"
)

func (m *Modulator) BuildWav(filename string) {
	randomNoiseStartDuration := 0 //(rand.Intn(5) + 2) * 1000
	randomNoiseEndDuration := 0   //(rand.Intn(5) + 2) * 1000

	rate := m.Options.Kilobitrate * 1000
	bitDepth := options.BitDepth
	//fmt.Println(randomNoiseStartDuration)

	randomNoiseStartBits := rate * (randomNoiseStartDuration / 1000)
	randomNoiseEndBits := rate * (randomNoiseEndDuration / 1000)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	buffer := audio.IntBuffer{}
	buffer.Format = &audio.Format{}
	buffer.Format.NumChannels = 1
	buffer.Format.SampleRate = rate

	encoder := aiff.NewEncoder(file, buffer.Format.SampleRate, bitDepth, buffer.Format.NumChannels)

	wave := m.FullWave()
	buffer.Data = make([]int, len(wave)+randomNoiseStartBits+randomNoiseEndBits)

	scaler := float64(math.Pow(2, float64(bitDepth-1)))

	offset := 0

	for i := 0; i < randomNoiseStartBits; i++ {
		noiseAmplitude := (scaler / float64(100.0)) * float64(m.Options.NoiseLevel)
		noise := noiseAmplitude * rand.Float64()

		buffer.Data[offset] = int(noise)
		offset++
	}

	for _, amplitude := range wave {
		buffer.Data[offset] = int(amplitude)
		offset++
	}

	for i := 0; i < randomNoiseEndBits; i++ {
		noiseAmplitude := (scaler / float64(100.0)) * float64(m.Options.NoiseLevel)
		noise := noiseAmplitude * rand.Float64()

		buffer.Data[offset] = int(noise)
		offset++
	}

	encoder.Write(&buffer)
	encoder.Close()
}
