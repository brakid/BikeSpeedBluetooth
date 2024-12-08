package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var upgrader = websocket.Upgrader{}
var trainer = Trainer{}
var bikeSensor = BikeSensor{2135}
var subscribers = CreateSubscribers()
var tracks = CreateTracks()
var wg = sync.WaitGroup{}
var indexTemplate = template.Must(template.New("template.html").ParseFiles("./template.html"))

func subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("reading subscriber id:", err)
		return
	}
	var subscriberId SubscriberId
	err = json.Unmarshal(message, &subscriberId)
	if err != nil {
		log.Println("parsing subscriber id:", err)
		return
	}

	defer c.Close()
	defer subscribers.Delete(c)
	subscribers.Add(c, subscriberId.ID)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			log.Printf("%v connection closed\n", c.RemoteAddr())
			break
		}
	}

	subscriber := subscribers.Get(c)
	if subscriber == nil {
		log.Println("no subscriber found")
		return
	}
	tracks.Add(subscriber.subscriberId, subscriber.GetTrainingData())
}

func getTrack(w http.ResponseWriter, r *http.Request) {
	trackId := r.URL.Query().Get("trackId")

	if trackId == "" {
		http.Error(w, "missing trackId", 400)
	}

	data, err := tracks.DownloadAsTCX(trackId)
	if (err != nil) || (len(data) == 0) {
		http.Error(w, "no data found", 404)
	}

	w.Header().Set("Content-Type", "application/vnd.garmin.tcx+xml")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=track-%v.tcx", trackId))
	w.Write([]byte(data))
}

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, r.Host)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/subscribe", subscribe)
		http.HandleFunc("/tracks", getTrack)
		http.HandleFunc("/", index)
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	err := ListenToSpeedAndCadenceSensor(adapter, "R2Duo", &bikeSensor, &trainer, subscribers)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	wg.Wait()
}
