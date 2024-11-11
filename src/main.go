package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var upgrader = websocket.Upgrader{}
var subscribers = make(map[*websocket.Conn]*Subscriber)
var wg = sync.WaitGroup{}
var homeTemplate = template.Must(template.New("template.html").ParseFiles("./template.html"))

func subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	defer delete(subscribers, c)
	subscribers[c] = CreateSubscriber(c)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.Execute(w, "ws://"+r.Host+"/subscribe")
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func main() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/subscribe", subscribe)
		http.HandleFunc("/", home)
		log.Printf("Error: %v\n", http.ListenAndServe("localhost:8080", nil))
	}()

	err := adapter.Enable()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Searching for R2Duo")
	var address *bluetooth.Address
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == "R2Duo" {
			log.Println("Found R2Duo")
			address = &device.Address
			adapter.StopScan()
		}
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if address == nil {
		log.Fatalf("No device found")
	}

	device, err := adapter.Connect(*address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDCyclingSpeedAndCadence})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	service := services[0]
	characteristics, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDCSCMeasurement})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Subscribing to speed and cadence events on R2Duo")
	char := characteristics[0]
	var lastData *SpeedCadenceData
	wheelRevolutionTimes, _ := CreateBuffer[*RevolutionTime](5)
	crankRevolutionTimes, _ := CreateBuffer[*RevolutionTime](5)

	char.EnableNotifications(func(buf []byte) {
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

					log.Printf("Wheel: %v kmh", GetWheelKmH(wheelRevolutions, wheelTimeDifference, 2135))
					log.Printf("Crank: %v rpM", GetCrankRpM(crankRevolutions, crankTimeDifference))

					for connection := range subscribers {
						subscribers[connection].Send(Speedometer{GetWheelKmH(wheelRevolutions, wheelTimeDifference, 2135), GetCrankRpM(crankRevolutions, crankTimeDifference)})
					}
				}
			}

			lastData = data
		}
	})

	wg.Wait()
}
