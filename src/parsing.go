package main

import (
	"fmt"
)

type SpeedCadenceData struct {
	WheelRevolutionsPresent    bool
	CrankRevolutionsPresent    bool
	CumulativeWheelRevolutions uint32
	CumulativeCrankRevolutions uint16
	LastWheelEventTime         uint16
	LastCrankEventTime         uint16
}

type RevolutionTime struct {
	Revolutions    uint32
	TimeDifference float64
}

func ParseData(data []byte) (*SpeedCadenceData, error) {
	if len(data) != 11 {
		return nil, fmt.Errorf("expected data to contain 11 bytes")
	}

	wheelRevolutionsPresent := (data[0] & 0b00000001) == 1
	crankRevolutionsPresent := ((data[0] & 0b00000010) >> 1) == 1
	cumulativeWheelRevolutions := uint32(0)
	cumulativeWheelRevolutions |= uint32(data[1])
	cumulativeWheelRevolutions |= uint32(data[2]) << 8
	cumulativeWheelRevolutions |= uint32(data[3]) << 16
	cumulativeWheelRevolutions |= uint32(data[4]) << 24
	lastWheelEventTime := uint16(0)
	lastWheelEventTime |= uint16(data[5])
	lastWheelEventTime |= uint16(data[6]) << 8
	cumulativeCrankRevolutions := uint16(0)
	cumulativeCrankRevolutions |= uint16(data[7])
	cumulativeCrankRevolutions |= uint16(data[8]) << 8
	lastCrankEventTime := uint16(0)
	lastCrankEventTime |= uint16(data[9])
	lastCrankEventTime |= uint16(data[10]) << 8

	speedCadenceData := SpeedCadenceData{
		WheelRevolutionsPresent:    wheelRevolutionsPresent,
		CrankRevolutionsPresent:    crankRevolutionsPresent,
		CumulativeWheelRevolutions: cumulativeWheelRevolutions,
		CumulativeCrankRevolutions: cumulativeCrankRevolutions,
		LastWheelEventTime:         lastWheelEventTime,
		LastCrankEventTime:         lastCrankEventTime,
	}

	return &speedCadenceData, nil
}

func CalculateUint32Difference(lastValue uint32, currentValue uint32) uint32 {
	enlargedLastValue := uint64(lastValue)
	enlargedCurrentValue := uint64(currentValue)
	if enlargedCurrentValue < enlargedLastValue {
		enlargedCurrentValue += uint64(1) << 32
	}
	return uint32(enlargedCurrentValue - enlargedLastValue)
}

func CalculateUint16Difference(lastValue uint16, currentValue uint16) uint16 {
	enlargedLastValue := uint32(lastValue)
	enlargedCurrentValue := uint32(currentValue)
	if enlargedCurrentValue < enlargedLastValue {
		enlargedCurrentValue += uint32(1) << 16
	}
	return uint16(enlargedCurrentValue - enlargedLastValue)
}

func CalculateDifference(lastSpeedCadenceData, currentSpeedCadenceData *SpeedCadenceData) (*RevolutionTime, *RevolutionTime, error) {
	if currentSpeedCadenceData == nil || lastSpeedCadenceData == nil {
		return nil, nil, fmt.Errorf("nil data passed")
	}

	wheelRevolutions := CalculateUint32Difference(lastSpeedCadenceData.CumulativeWheelRevolutions, currentSpeedCadenceData.CumulativeWheelRevolutions)
	wheelTimeDifference := float64(CalculateUint16Difference(lastSpeedCadenceData.LastWheelEventTime, currentSpeedCadenceData.LastWheelEventTime)) / 1024
	crankRevolutions := uint32(CalculateUint16Difference(lastSpeedCadenceData.CumulativeCrankRevolutions, currentSpeedCadenceData.CumulativeCrankRevolutions))
	crankTimeDifference := float64(CalculateUint16Difference(lastSpeedCadenceData.LastCrankEventTime, currentSpeedCadenceData.LastCrankEventTime)) / 1024

	wheelRevolutionTime := RevolutionTime{
		Revolutions:    wheelRevolutions,
		TimeDifference: wheelTimeDifference,
	}
	crankRevolutionTime := RevolutionTime{
		Revolutions:    crankRevolutions,
		TimeDifference: crankTimeDifference,
	}
	return &wheelRevolutionTime, &crankRevolutionTime, nil
}

func Sum(revolutionTimes []*RevolutionTime) (uint64, float64) {
	totalRevolutions := uint64(0)
	totalTimeDifference := float64(0)
	for _, revolutionTime := range revolutionTimes {
		if revolutionTime == nil {
			continue
		}
		totalRevolutions += uint64(revolutionTime.Revolutions)
		totalTimeDifference += revolutionTime.TimeDifference
	}

	return totalRevolutions, totalTimeDifference
}

func ConvertToDistanceInKm(revolutions uint64, distanceInMm float64) float64 {
	return float64(revolutions) * (distanceInMm / 1_000_000)
}

func GetWheelKmH(revolutions uint64, timeDifferenceInS float64, distanceInMm float64) float64 {
	if timeDifferenceInS == 0 {
		return 0.0
	}
	return 3600 * ConvertToDistanceInKm(revolutions, distanceInMm) / timeDifferenceInS
}

func GetCrankRpM(revolutions uint64, timeDifferenceInS float64) float64 {
	if timeDifferenceInS == 0 {
		return 0.0
	}
	return 60 * (float64(revolutions)) / (timeDifferenceInS)
}
