package demodulation

import (
	"github.com/Hunter-Dolan/midrange/fec"
)

type frame struct {
	index                int
	demodulator          *Demodulator
	confidenceCollection *fec.DataConfidenceCollection
}

func (f *frame) demodulate() {
	powers := f.demodulator.carrierPowerAtFrame(f.index)

	trainer := f.demodulator.trainer

	f.confidenceCollection = &fec.DataConfidenceCollection{}

	for carrierIndex, power := range powers {
		value, confidence := trainer.determineValueAndConfidence(power, carrierIndex)
		f.confidenceCollection.Add(value, confidence)
	}

	//fmt.Println(f.confidenceCollection.Data())

}
