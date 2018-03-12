package domain

import (
	"time"
)

var (
	// ALLOWED_CABS has information about different cabs and their associated cab types that
	// application right now caters to.
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

// CabService is a serive which is an interface which exposes the method EtaNow which takes a pointer
// to a CbRequest as input and return a time.Duration and error as output.
// This tells you what is the eta for the request to that particular cab service.
// Any cab service (be it ola, uber, lyft) can implement the EtaNow method.
type CabService interface {
	EtaNow(*CabRequest) (time.Duration, error)
}

// CabBookingResponse is the final response the is generated of the application which is
// sent to the user as notification at the user's notificatio address,
// CabBookingResponse encapsulated information like the UserRequest and the best booking time
// to request/book a cab.
type CabBookingResponse struct {
	BookingID uint64
	*UserRequest
	BestBookingTime time.Time
}

// CabRequest is a composition of the attricutes which make a valid request
// which can be sent to the CabService.
type CabRequest struct {
	Source      Location
	Destination Location
	BookingTime time.Time
	Cab         string
	CabType     string
}

// CabBookingResponseRepository exposes the interface to store and find the cab booking responses
// from a repository.
type CabBookingResponseRepository interface {
	FindById(uint64) *CabBookingResponse
	Store(*CabBookingResponse) (uint64, error)
}

// validateCab if a function which takes cab and cabType, both of type string as inputs and returns
// if the cab, cabType combination is valid and allowed by the applicaion or not. It returns a bool.
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

// NewCabRequest is a constructor that takes in many attributes which form the CabRequest object and
// constructs a new CabRequst object with those inputs and send a pointer to that objec in return.
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

// NewCabBookingResponse is a constructor function which takes pointers to UserRequest object and the
// bestBookingTime of type time.Time as inputs and returns a pointer to the newly created
// CabBookingResponse object.
func NewCabBookingResponse(ur *UserRequest, bestBookingTime time.Time) *CabBookingResponse {
	cr := CabBookingResponse{
		UserRequest:     ur,
		BestBookingTime: bestBookingTime,
	}

	return &cr

}
