package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02T15:04:05.000Z"

type Tracks struct {
	tracks map[string][]Speedometer
}

func CreateTracks() *Tracks {
	return &Tracks{tracks: make(map[string][]Speedometer)}
}

func (t *Tracks) Add(trackId string, track []Speedometer) {
	t.tracks[trackId] = track
}

func (t *Tracks) DownloadAsTCX(trackId string) (string, error) {
	track, isFound := t.tracks[trackId]
	if !isFound {
		return "", fmt.Errorf("no track with id %v found", trackId)
	}
	return GetTrackTCX(track), nil
}

func GetTrackpoint(data Speedometer) string {
	timestamp := time.Unix(0, int64(data.Timestamp)*int64(time.Millisecond)).UTC()

	return fmt.Sprintf(`
<Trackpoint>
	<Time>%v</Time>
	<DistanceMeters>%.6f</DistanceMeters>
	<Cadence>%d</Cadence>
	<Extensions>
		<TPX xmlns="http://www.garmin.com/xmlschemas/ActivityExtension/v2" CadenceSensor="Bike">
			<Speed>%.6f</Speed>
			<Watts>%d</Watts>
		</TPX>
	</Extensions>
	<HeartRateBpm xsi:type="HeartRateInBeatsPerMinute_t">
		<Value>0</Value>
	</HeartRateBpm>
</Trackpoint>
`, timestamp.Format(TIME_FORMAT), data.Distance*1000.0, uint64(math.Floor(data.Cadence)), data.Speed, uint64(math.Floor(data.Power)))
}

func GetTrackTCX(track []Speedometer) string {
	timestamp := time.Unix(0, int64(track[0].Timestamp)*int64(time.Millisecond)).UTC()

	var builder strings.Builder
	for _, trackpoint := range track {
		builder.WriteString(GetTrackpoint(trackpoint))
	}

	totalSeconds := uint64(math.Floor(track[len(track)-1].Duration))
	totalDistance := uint64(math.Floor(track[len(track)-1].Distance * 1000))

	return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<TrainingCenterDatabase xmlns="http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.garmin.com/xmlschemas/ActivityExtension/v2 http://www.garmin.com/xmlschemas/ActivityExtensionv2.xsd http://www.garmin.com/xmlschemas/FatCalories/v1 http://www.garmin.com/xmlschemas/fatcalorieextensionv1.xsd http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd">
	<Activities>
		<Activity Sport="Biking">
			<Id>%v</Id>
			<Lap StartTime="%v">
				<TotalTimeSeconds>%d</TotalTimeSeconds>
				<DistanceMeters>%d</DistanceMeters>
				<MaximumSpeed>%d</MaximumSpeed>
				<AverageHeartRateBpm xsi:type="HeartRateInBeatsPerMinute_t">
					<Value>0</Value>
				</AverageHeartRateBpm>
				<MaximumHeartRateBpm xsi:type="HeartRateInBeatsPerMinute_t">
					<Value>0</Value>
				</MaximumHeartRateBpm>
				<Intensity>Resting</Intensity>
				<TriggerMethod>Distance</TriggerMethod>
				<Track>
					%v
				</Track>
			</Lap>
		</Activity>
	</Activities>
</TrainingCenterDatabase>
`, timestamp.Format(TIME_FORMAT), timestamp.Format(TIME_FORMAT), totalSeconds, totalDistance, FindMaxSpeed(track), builder.String())
}
