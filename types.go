package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

func CreateSubscriber(conn *websocket.Conn) *Subscriber {
	subscriber := Subscriber{conn, sync.Mutex{}}
	return &subscriber
}

func (s *Subscriber) Send(data interface{}) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	return s.conn.WriteJSON(data)
}

type Speedometer struct {
	Speed   float64 `json:"speed"`
	Cadence float64 `json:"cadence"`
}
