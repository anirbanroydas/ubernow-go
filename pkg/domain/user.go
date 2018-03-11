// package domain consists of all the domain object the business
// ubernow-go can be considered in general as a service which takes requests
// from users and notifiies them when is the right time to book a cab,
// now, this cab can be uber, ola, etc, for simnplicity we are talking about uber,
// but the business considers it just a cab.
//
// So using DDD principles I have found the domain objects
// Lets talk a little about the domain experts, this may not be exactly a direct
// application of DDD but a flavour which is easier to implement ond practical to
// proceed with, at a faster pace (to keep the agile going)
//
// Domain experts:
// 1. User - user knows about their requests, like source, destination, reaching time,
// their address to get the notifications at, the cab and the cab type they want to use
// 2. Cab Service - cab service is the domain expert on everything related to cabs,
// how to send a booking request, what is the source and destination of the booking,
// what is the cab availability, what is the eta for a booked cab, what is limitation of
// booking, who is a user for their service, etc
// 3. Traffic Service - traffice service is the domian expert on traffic related data,
// like sourc and destination of a traffic route, the traffic condition, the time of the
// day for the traffic data, the best case travel time, the worst case travel time
// between the source and destination given the time of the day and the current
// traffic condition(like busy, jam, free, etc)
//
// So by using the **Domain Experts** knowledge and seeing the **Ubiquitous Language**,
// we can derive at some domain objects
//
// - Value Object - Things whose different instances having same values are identical
// - Entities - Things whose different instance having same values are different
// - Aggregates - A collection of Entities and Value Objects which are having high cohesion
// and it also has a root entinty which is the main source of communicting with the other
// entities in the aggregate, also without the root the othe entities doesn't make sense
// or will never be used directly, in the business use case
// - Services - Things which are not a single entity or aggregrage but provide some kind
// of behaviour or service by using the aggregrates, entities etc.
// - Repository - Thigns which facilitate the access and storage of entities and
// aggregrates (if required).
//
// This package basically creates all the domain object and expose some interfaces
// to use them but doesnt depend on anything outside the domain layer
package domain

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	reaching_time_threshold_in_minute int
	default_time_threshold_in_minute  = 5
)

// user is an  entity which encapsulate information aobut a user.
type User struct {
	UserID uint64
	Name   string
}

// location encapsulates any kind of location like source, destination used
// by other domain objects.
type Location struct {
	Name      string
	Latitude  string
	Longitude string
}

// userAddress encapsulated the different addresses that belong to the user, like
// where the user wants to get notified at, so the addrType can be lets say "email"
// and the value can be "anirban.nick@gmail.com", or it can be addrType: "sms"
// and value: "xxxxxxxxxx".
type UserAddress struct {
	AddrType string
	Value    string
}

// request is the root of the aggregrae which encapsulate information like source,
// destination, reaching time of the user, cab and cab type preferred, notification address
// of user where the user needs the notification regrading when to book the cab is to be sent.
type Request struct {
	reqID            uint64
	Source           Location
	Destination      Location
	ReachingTime     time.Time
	Cab              string
	CabType          string
	NotificationAddr UserAddress
}

// userRequest encapsulates a request with a user
type UserRequest struct {
	*User
	*Request
}

// UserRepository exposes the interface to store and find user from a repository.
type UserRepository interface {
	FindByID(uint64) (*User, error)
	Store(*User) (uint64, error)
}

// RequestRepository exposes the interface to store and find requests from a repository.
type RequestRepository interface {
	FindByID(uint64) (*Request, error)
	Store(*Request) (uint64, error)
}

type UserAddressValidator interface {
	Validate(UserAddress) error
}

type Logger interface {
	LogError(string)
	LogInfo(string)
}

// type LocationValidator interface {
// 	Validate(location) bool
// }

// type TimeValidator interface {
// 	Validate(time.Time) error
// }

// type CabValidator interface {
// 	Validate(cab, cabType string) bool
// }

func init() {
	reaching_time_threshold_string, ok := os.LookupEnv("REACHING_TIME_THRESHOLD")
	if !ok {
		reaching_time_threshold_in_minute = default_time_threshold_in_minute
		return
	}
	var err error
	reaching_time_threshold_in_minute, err = strconv.Atoi(reaching_time_threshold_string)
	if err != nil {
		reaching_time_threshold_in_minute = default_time_threshold_in_minute
	}
}

func ToString(t time.Time) string {
	h, m, s := "", "", ""
	if t.Hour() != 0 {
		h = strconv.Itoa(t.Hour()) + "h"
	}
	if t.Minute() != 0 {
		m = strconv.Itoa(t.Minute()) + "m"
	}
	if t.Second() != 0 {
		s = strconv.Itoa(t.Second()) + "s"
	}
	return h + m + s
}

func validateLocation(l Location) bool {
	if l.Latitude == "" || l.Longitude == "" {
		return false
	}
	return true
}

func validateReachingTime(rt time.Time) error {
	if rt.Sub(time.Now()) < time.Duration(time.Duration(reaching_time_threshold_in_minute)*time.Minute) {
		return errors.New(fmt.Sprintf("reaching time: %s has past or is less than threshold interval: %d minutes from current time", rt, reaching_time_threshold_in_minute))
	}
	return nil
}

func NewRequest(source, destination Location, reachingTime time.Time, cab, cabType string, notificationAddr UserAddress, uav UserAddressValidator) (*Request, error) {
	var r *Request
	ok := validateLocation(source)
	if !ok {
		return r, errors.New(fmt.Sprintf("source location: %v is not valid", source))
	}
	ok = validateLocation(destination)
	if !ok {
		return r, errors.New(fmt.Sprintf("destination location: %v is not valid", destination))
	}
	err := validateReachingTime(reachingTime)
	if err != nil {
		return r, errors.Wrap(err, "NewRequest failed for timeValidator error")
	}
	ok = validateCab(cab, cabType)
	if !ok {
		return r, errors.New(fmt.Sprintf("requested cab: %s or cabtype: %s not avaialable", cab, cabType))
	}
	err = uav.Validate(notificationAddr)
	if err != nil {
		return r, errors.Wrap(err, "NewRequest couldn't validate notification address")
	}

	r = &Request{
		Source:           source,
		Destination:      destination,
		ReachingTime:     reachingTime,
		Cab:              cab,
		CabType:          cabType,
		NotificationAddr: notificationAddr,
	}

	return r, nil
}

func NewUser(name string) *User {
	u := User{
		Name: name,
	}
	return &u
}

func NewUserRequest(u *User, r *Request) *UserRequest {
	ur := UserRequest{
		User:    u,
		Request: r,
	}

	return &ur
}
