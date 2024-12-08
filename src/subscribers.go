package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	conn           *websocket.Conn
	subscriberId   string
	mutex          sync.Mutex
	distanceOffset float64
	durationOffset float64
	initialized    bool
	data           []Speedometer
}

func CreateSubscriber(conn *websocket.Conn, subscriberId string) *Subscriber {
	subscriber := Subscriber{conn, subscriberId, sync.Mutex{}, 0, 0, false, make([]Speedometer, 0)}
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

	speedometerData := Speedometer{data.Timestamp, data.Speed, data.Cadence, data.Power, data.Distance - s.distanceOffset, data.Duration - s.durationOffset}
	s.data = append(s.data, speedometerData)
	return s.conn.WriteJSON(&speedometerData)
}

func (s *Subscriber) GetTrainingData() []Speedometer {
	return s.data
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

func (s *Subscribers) Add(conn *websocket.Conn, subscriberId string) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.subscribers[conn] = CreateSubscriber(conn, subscriberId)
}

func (s *Subscribers) Delete(conn *websocket.Conn) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	delete(s.subscribers, conn)
}

func (s *Subscribers) Get(conn *websocket.Conn) *Subscriber {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	subscriber := s.subscribers[conn]
	return subscriber
}

func (s *Subscribers) Send(data *Speedometer) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	for connection := range s.subscribers {
		s.subscribers[connection].Send(data)
	}
}
