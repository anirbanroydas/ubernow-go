package domain

import (
	"time"
)

// TrafficService is a service which is an interface which exposed the method TravelTime which takes in
// a pointer to TrafficRequest object as the input and returns a pointer to TrafficResponse object along
// with an error.
// And traffie Serive like googlemaps, mapbox, etc. can be used to implement this service.
type TrafficService interface {
	TravelTime(*TrafficRequest) (*TrafficResponse, error)
}

// TrafficRequest is the encapsulation of the data needed to create a valid request
// that can be sent to the TrafficService.
type TrafficRequest struct {
	Source      Location
	Destination Location
	TimeOfDay   time.Time
}

// TrafficResponse is the response sent by TrafficService for a corresponding TrafficRequest. The object also
// ecnapsulates a pointer to the TrafficRequest.
type TrafficResponse struct {
	*TrafficRequest
	TravelTime time.Time
	BestCase   time.Time
	WorstCase  time.Time
}

// NewTrafficRequest is a constructor function which takes Location and time.Time attribute necessary to
// construct a new TrafficRequest and returns a pointer to that object.
func NewTrafficRequest(source, destination Location, timeOfDay time.Time) *TrafficRequest {
	tr := TrafficRequest{
		Source:      source,
		Destination: destination,
		TimeOfDay:   timeOfDay,
	}

	return &tr
}
