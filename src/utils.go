package main

import "math"

func FindMaxSpeed(track []Speedometer) uint64 {
	maxSpeed := uint64(0)
	for _, data := range track {
		speed := uint64(math.Floor(data.Speed))
		if speed > maxSpeed {
			maxSpeed = speed
		}
	}

	return maxSpeed
}
