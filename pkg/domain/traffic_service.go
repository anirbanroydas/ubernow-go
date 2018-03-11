package domain

import (
	"time"
)

// trafficRequest is the encapsulation of the data needed to request the traffic service.
type TrafficRequest struct {
	Source      Location
	Destination Location
	TimeOfDay   time.Time
}

// trafficResponse is the result the traffic service send for the corresponding trafficRequest.
type TrafficResponse struct {
	*TrafficRequest
	TravelTime time.Time
	BestCase   time.Time
	WorstCase  time.Time
}

// trafficService is a service which exposes the interface having the method travelTime which takes in
// trafficREquest as the input and returns a trafficResponse. And traffieSerive like googlemaps, mapbox,
// etc can be used to implement this service.
type TrafficService interface {
	TravelTime(*TrafficRequest) (*TrafficResponse, error)
}

func NewTrafficRequest(source, destination Location, timeOfDay time.Time) *TrafficRequest {
	tr := TrafficRequest{
		Source:      source,
		Destination: destination,
		TimeOfDay:   timeOfDay,
	}

	return &tr
}
