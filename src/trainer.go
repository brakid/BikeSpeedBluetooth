package main

import "math"

type Trainer struct{}

func (t *Trainer) CalculatePower(speedKmh float64, cadenceRpm float64) float64 {
	if cadenceRpm == 0 {
		return 0.0
	}

	power := -0.00009*math.Pow(speedKmh, 3) + 0.1173*math.Pow(speedKmh, 2) + 3.437*speedKmh - 1.236

	if power < 0 {
		return 0.0
	}

	return power
}
