package main

import (
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
)

func ListenToSpeedAndCadenceSensor(adapter *bluetooth.Adapter, deviceName string, bikeSensor *BikeSensor, trainer *Trainer, subscribers *Subscribers) error {
	if adapter == nil || subscribers == nil {
		return fmt.Errorf("nil adapter or subscribers")
	}

	err := adapter.Enable()
	if err != nil {
		return err
	}

	log.Printf("Searching for %v\n", deviceName)
	var address *bluetooth.Address
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == deviceName {
			log.Printf("Found %v\n", deviceName)
			address = &device.Address
			adapter.StopScan()
		}
	})
	if err != nil {
		return err
	}
	if address == nil {
		return fmt.Errorf("no device found")
	}

	device, err := adapter.Connect(*address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDCyclingSpeedAndCadence})
	if err != nil {
		return err
	}

	service := services[0]
	characteristics, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDCSCMeasurement})
	if err != nil {
		return err
	}
	log.Printf("Subscribing to speed and cadence events on %v\n", deviceName)
	characteristic := characteristics[0]
	var lastData *SpeedCadenceData
	wheelRevolutionTimes, _ := CreateBuffer[*RevolutionTime](5)
	crankRevolutionTimes, _ := CreateBuffer[*RevolutionTime](5)
	totalWheelRevolutions := uint64(0)
	totalDuration := float64(0)

	characteristic.EnableNotifications(func(buf []byte) {
		data, err := ParseData(buf)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			if data != nil && lastData != nil {
				wheelRevolutionTime, crankRevolutionTime, err := CalculateDifference(lastData, data)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					wheelRevolutionTimes.Append(wheelRevolutionTime)
					crankRevolutionTimes.Append(crankRevolutionTime)

					wheelRevolutions, wheelTimeDifference := Sum(wheelRevolutionTimes.Get())
					crankRevolutions, crankTimeDifference := Sum(crankRevolutionTimes.Get())
					totalWheelRevolutions += uint64(wheelRevolutionTime.Revolutions)
					totalDuration += wheelRevolutionTime.TimeDifference

					wheelKmh := bikeSensor.CalculateWheelKmH(wheelRevolutions, wheelTimeDifference)
					crankRpm := CalculateCrankRpM(crankRevolutions, crankTimeDifference)
					power := trainer.CalculatePower(wheelKmh, crankRpm)
					distance := bikeSensor.CalculateDistanceKmH(totalWheelRevolutions)

					log.Printf("Wheel: %v kmh", wheelKmh)
					log.Printf("Crank: %v rpM", crankRpm)
					log.Printf("Power: %v watts", power)
					log.Printf("Distance: %v km", distance)
					log.Printf("Duration: %v s", totalDuration)

					subscribers.Send(
						&Speedometer{
							wheelKmh,
							crankRpm,
							power,
							distance,
							totalDuration})
				}
			}

			lastData = data
		}
	})

	return nil
}
