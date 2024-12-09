package main

import "math"

type Trainer struct{}

func (t *Trainer) CalculatePower(speedKmh float64, cadenceRpm float64) float64 {
	if cadenceRpm == 0 {
		return 0.0
	}

	power := 0.00004667*math.Pow(speedKmh, 4) - 0.005437*math.Pow(speedKmh, 3) + 0.2577*math.Pow(speedKmh, 2) + 2.050*speedKmh

	if power < 0 {
		return 0.0
	}

	return power
}
