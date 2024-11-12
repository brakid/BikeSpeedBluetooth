package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	conn           *websocket.Conn
	mutex          sync.Mutex
	distanceOffset float64
	durationOffset float64
	initialized    bool
}

func CreateSubscriber(conn *websocket.Conn) *Subscriber {
	subscriber := Subscriber{conn, sync.Mutex{}, 0, 0, false}
	return &subscriber
}

func (s *Subscriber) Send(data *Speedometer) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if !s.initialized {
		s.distanceOffset = data.Distance
		s.durationOffset = data.Duration
		s.initialized = true
	}
	return s.conn.WriteJSON(&Speedometer{data.Speed, data.Cadence, data.Distance - s.distanceOffset, data.Duration - s.durationOffset})
}

type Subscribers struct {
	subscribers map[*websocket.Conn]*Subscriber
	mutex       sync.Mutex
}

func CreateSubscribers() *Subscribers {
	return &Subscribers{
		make(map[*websocket.Conn]*Subscriber),
		sync.Mutex{},
	}
}

func (s *Subscribers) Add(conn *websocket.Conn) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.subscribers[conn] = CreateSubscriber(conn)
}

func (s *Subscribers) Delete(conn *websocket.Conn) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	delete(s.subscribers, conn)
}

func (s *Subscribers) Send(data *Speedometer) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	for connection := range s.subscribers {
		s.subscribers[connection].Send(data)
	}
}
