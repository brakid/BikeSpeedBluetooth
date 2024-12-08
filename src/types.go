package main

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

type SubscriberId struct {
	ID string `json:"id"`
}

type Speedometer struct {
	Timestamp uint64  `json:"timestamp"` // datum timestamp as Unix Epoch Millis
	Speed     float64 `json:"speed"`     // speed in kmH
	Cadence   float64 `json:"cadence"`   // cadence in RpM
	Power     float64 `json:"power"`     // power in watts
	Distance  float64 `json:"distance"`  // ditance in km
	Duration  float64 `json:"duration"`  // dureaction in seconds since started
}
