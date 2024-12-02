package main

type Speedometer struct {
	Speed    float64 `json:"speed"`
	Cadence  float64 `json:"cadence"`
	Power  float64 `json:"power"`
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}
