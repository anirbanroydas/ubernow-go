package domain

import (
	"fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		name     string
		tobj     time.Time
		expected string
	}{
		{
			name:     "time with hour month second",
			tobj:     time.Date(2018, time.November, 3, 5, 6, 7, 0, time.UTC),
			expected: "5h6m7s",
		},
		{
			name:     "time with hour",
			tobj:     time.Date(2018, time.November, 3, 8, 0, 0, 0, time.UTC),
			expected: "8h",
		},
		{
			name:     "time with month second",
			tobj:     time.Date(2018, time.November, 3, 0, 7, 6, 0, time.UTC),
			expected: "7m6s",
		},
		{
			name:     "time with second",
			tobj:     time.Date(2018, time.November, 3, 0, 0, 4, 0, time.UTC),
			expected: "4s",
		},
		{
			name:     "time with hour month",
			tobj:     time.Date(2018, time.November, 3, 8, 9, 0, 0, time.UTC),
			expected: "8h9m",
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result := ToString(tc.tobj)
			if result != tc.expected {
				t.Errorf("%s: ToString(%s) => got: %s, expected: %s", tc.name, tc.tobj, result, tc.expected)
			}
		})
	}
}

func TestValidateLocation(t *testing.T) {
	testCases := []struct {
		name     string
		loc      Location
		expected bool
	}{
		{
			name: "valid location with name",
			loc: Location{
				Name:      "koramangala",
				Latitude:  "77.12341324",
				Longitude: "23.1341324312",
			},
			expected: true,
		},
		{
			name: "valid location without name",
			loc: Location{
				Latitude:  "77.12341324",
				Longitude: "23.1341324312",
			},
			expected: true,
		},
		{
			name:     "empty location",
			loc:      Location{},
			expected: false,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result := validateLocation(tc.loc)
			if result != tc.expected {
				t.Errorf("%s: validateLocation(%v) => Got: %v, expected: %v", tc.name, tc.loc, result, tc.expected)
			}
		})
	}

}

func TestValidateReachingTime(t *testing.T) {
	rt1 := time.Now().Add(100 * time.Minute)
	rt2 := time.Now().Add(-(1 * time.Minute))
	rt3 := time.Now().Add(3 * time.Minute)

	testCases := []struct {
		name         string
		reachingTime time.Time
		expected     error
	}{
		{
			name:         "valid reaching time - time after time.Now() - after 5 minitues of time.Now()",
			reachingTime: rt1,
			expected:     nil,
		},
		{
			name:         "invalid reaching time - time before time.Now()",
			reachingTime: rt2,
			expected:     errors.New(fmt.Sprintf("reaching time: %s has past or is less than threshold interval: %d minutes from current time", rt2, 5)),
		},
		{
			name:         "invalid reaching time - time after time.Now() - within 5 mintues of time.Now()",
			reachingTime: rt3,
			expected:     errors.New(fmt.Sprintf("reaching time: %s has past or is less than threshold interval: %d minutes from current time", rt3, 5)),
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result := validateReachingTime(tc.reachingTime)
			if (result != nil && tc.expected == nil) || (result == nil && tc.expected != nil) {
				t.Errorf("%s: validateReachingTime(%v) => Got: %v, expected: %v", tc.name, tc.reachingTime, result, tc.expected)
			}
		})
	}
}
