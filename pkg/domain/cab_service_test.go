package domain

import (
	// "fmt"
	// "reflect"
	"testing"
)

func TestValidateCab(t *testing.T) {
	testCases := []struct {
		name         string
		cab, cabType string
		expected     bool
	}{
		{
			name: "valid cab, valid cabType",
			cab:  "uber", cabType: "uberGo",
			expected: true,
		},
		{
			name: "valid cab, invalid cabType",
			cab:  "uber", cabType: "uberMe",
			expected: false,
		},
		{
			name: "invalid cab, valid cabType",
			cab:  "lyft", cabType: "uberGo",
			expected: false,
		},
		{
			name: "invalid cab, invalid cabType",
			cab:  "lyft", cabType: "micro",
			expected: false,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result := validateCab(tc.cab, tc.cabType)
			if result != tc.expected {
				t.Errorf("%s: validateCab(%s, %s) => Got: %v, expected: %v", tc.name, tc.cab, tc.cabType, result, tc.expected)
			}
		})
	}
}
