package transmission

import "github.com/Hunter-Dolan/midrange/demodulation"

func (t *Transmission) SetWave64(wave []float64) {
	t.wave = &wave
	t.demodulate()
}

func (t *Transmission) SetWave(wave []float32) {

	float64Wave := make([]float64, len(wave))

	for i, amp := range wave {
		float64Wave[i] = float64(amp)
	}

	t.wave = &float64Wave
	t.demodulate()
}

func (t *Transmission) Data() []byte {
	return t.demodulator.Data()
}

func (t *Transmission) demodulate() {
	if t.demodulator == nil {
		t.demodulator = &demodulation.Demodulator{}
		t.demodulator.Options = t.options

		start := 0
		end := len(*t.wave)

		wave := (*t.wave)[start:end]
		t.demodulator.SetWave(&wave)
	}
}
