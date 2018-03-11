package domain

import (
	"time"
)

var (
	ALLOWED_CABS = []struct {
		Name  string
		Types []string
	}{
		{
			Name: "uber",
			Types: []string{
				"uberGo",
				"uberBlack",
				"uberShare",
				"uberX",
			},
		},
	}
)

// cabService is a serive which exposes the interface having the method etaNow which takes a cabRequest
// as input and return a time.Time which tells what is the eta for the request. Any cab service be it
// ola, uber, can implement it.
type CabService interface {
	EtaNow(*CabRequest) (time.Duration, error)
}

// booking is response sent by the business to the user at the user's notificatio address,
// booking encapsulated information like the userRequest, booking time and eta of the cab.
type CabBookingResponse struct {
	BookingID uint64
	*UserRequest
	BestBookingTime time.Time
}

// cabRequest is a composition of request type which is sent to the cab service.
type CabRequest struct {
	Source      Location
	Destination Location
	BookingTime time.Time
	Cab         string
	CabType     string
}

type CabBookingResponseRepository interface {
	FindById(uint64) *CabBookingResponse
	Store(*CabBookingResponse) (uint64, error)
}

func validateCab(cab, cabType string) bool {
	for _, c := range ALLOWED_CABS {
		if cab == c.Name {
			for _, ct := range c.Types {
				if cabType == ct {
					return true
				}
			}
		}
	}

	return false
}

func NewCabRequest(source, destination Location, bookingTime time.Time, cab, cabType string) *CabRequest {
	cr := CabRequest{
		Source:      source,
		Destination: destination,
		BookingTime: bookingTime,
		Cab:         cab,
		CabType:     cabType,
	}

	return &cr
}

func NewCabBookingResponse(ur *UserRequest, bestBookingTime time.Time) *CabBookingResponse {
	cr := CabBookingResponse{
		UserRequest:     ur,
		BestBookingTime: bestBookingTime,
	}

	return &cr

}
