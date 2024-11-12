package main

type Speedometer struct {
	Speed    float64 `json:"speed"`
	Cadence  float64 `json:"cadence"`
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}
