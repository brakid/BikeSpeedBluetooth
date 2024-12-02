package main

type BikeSensor struct {
	wheelCircumferenceInMm float64
}

func (b *BikeSensor) CalculateWheelKmH(revolutions uint64, timeDifferenceInS float64) float64 {
	if timeDifferenceInS == 0 {
		return 0.0
	}
	return 3600 * b.CalculateDistanceKmH(revolutions) / timeDifferenceInS
}

func (b *BikeSensor) CalculateDistanceKmH(revolutions uint64) float64 {
	return float64(revolutions) * (b.wheelCircumferenceInMm / 1_000_000)
}
