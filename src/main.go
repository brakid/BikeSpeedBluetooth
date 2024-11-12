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
var subscribers = CreateSubscribers()
var wg = sync.WaitGroup{}
var indexTemplate = template.Must(template.New("template.html").ParseFiles("./template.html"))

func subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	defer subscribers.Delete(c)
	subscribers.Add(c)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, "ws://"+r.Host+"/subscribe")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/subscribe", subscribe)
		http.HandleFunc("/", index)
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	err := ListenToSpeedAndCadenceSensor(adapter, "R2Duo", 2135, subscribers)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	wg.Wait()
}
