package demodulation

import "math"

type trainer struct {
	demodulator *Demodulator

	carrierActivePowers   []float64
	carrierInactivePowers []float64

	carrierActivePowerAverage   float64
	carrierInactivePowerAverage float64
}

func (t *trainer) train() {
	evenFrame := t.demodulator.carrierPowerAtFrame(0)
	oddFrame := t.demodulator.carrierPowerAtFrame(1)

	t.carrierActivePowers = make([]float64, t.demodulator.carrierCount)
	t.carrierInactivePowers = make([]float64, t.demodulator.carrierCount)

	activeSum := 0.0
	inactiveSum := 0.0

	for i, evenFrameValue := range evenFrame {
		oddFrameValue := oddFrame[i]

		carrierEven := i % 2

		activePower := evenFrameValue
		inactivePower := oddFrameValue

		if carrierEven == 1 {
			activePower = oddFrameValue
			inactivePower = evenFrameValue
		}

		t.carrierActivePowers[i] = activePower
		t.carrierInactivePowers[i] = inactivePower

		activeSum += activePower
		inactiveSum += inactivePower
	}

	t.carrierActivePowerAverage = activeSum / float64(t.demodulator.carrierCount)
	t.carrierInactivePowerAverage = inactiveSum / float64(t.demodulator.carrierCount)
}

func (t *trainer) determineValueAndConfidence(power float64, carrierIndex int) (bool, float64) {
	activeDistance := math.Abs(power - t.carrierActivePowerAverage)
	inactiveDistance := math.Abs(power - t.carrierInactivePowerAverage)

	activeCarrierPower := t.carrierActivePowers[carrierIndex]
	inactiveCarrierPower := t.carrierInactivePowers[carrierIndex]

	activeCarrierDistance := math.Abs(power - activeCarrierPower)
	inactiveCarrierDistance := math.Abs(power - inactiveCarrierPower)

	value := activeDistance < inactiveDistance
	confidence := math.Abs((activeDistance + activeCarrierDistance) - (inactiveDistance + inactiveCarrierDistance))

	return value, confidence
}
