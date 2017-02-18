package transmission

import "github.com/Hunter-Dolan/midrange/demodulation"

func (t *Transmission) SetWave(wave []float64) {
	t.wave = &wave
	t.demodulate()
}

func (t *Transmission) Data() []byte {
	return t.demodulator.Data()
}

func (t *Transmission) demodulate() {
	if t.demodulator == nil {
		t.demodulator = &demodulation.Demodulator{}
		t.demodulator.Options = t.options
		t.demodulator.SetWave(t.wave)
	}
}
