package transmission

import (
	"github.com/Hunter-Dolan/midrange/demodulation"
	"github.com/Hunter-Dolan/midrange/modulation"
	"github.com/Hunter-Dolan/midrange/options"
)

type Transmission struct {
	header     *packet
	data       []*packet
	dataLength int

	originalDataLength int

	wave *[]float64

	modulator   *modulation.Modulator
	demodulator *demodulation.Demodulator
	options     *options.Options
}

func NewTransmission() *Transmission {
	t := Transmission{}

	return &t
}

func (t *Transmission) SetOptions(options options.Options) {
	t.options = &options
}
