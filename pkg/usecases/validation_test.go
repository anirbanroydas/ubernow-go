package usecases

import (
	// "fmt"
	// "reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type MockAddressValidator struct{}

func (d MockAddressValidator) Validate(ua domain.UserAddress) error {
	return nil
}

type MockBadAddressValidator struct{}

func (d MockBadAddressValidator) Validate(ua domain.UserAddress) error {
	return errors.New("some error")
}

func TestValidEmailValue(t *testing.T) {
	testCases := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "anirban.nick@gmail.com",
			expected: true,
		},
		{
			name:     "valid email with special characters",
			email:    "anirban_12.34.@gmail.co",
			expected: false,
		},
		{
			name:     "invalid email without a single @",
			email:    "anirba.nick[at]gmail.com",
			expected: false,
		},
		{
			name:     "invalid email without dot(.) after @",
			email:    "anirban.nick@gmail",
			expected: false,
		},
		{
			name:     "invalid email, without @ and dot(.) after @",
			email:    "anirba.nick",
			expected: false,
		},
		{
			name:     "invalid email with special, illegal characters",
			email:    "anirban@nic$2#me@gmail.com",
			expected: false,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			result := validEmailValue(tc.email)
			if result != tc.expected {
				t.Errorf("%s: validEmailValue(%s) => Got: %v, expected: %v", tc.name, tc.email, result, tc.expected)
			}
		})
	}
}

func TestNewUserAddressValidator(t *testing.T) {
	var nilAddressValidator domain.UserAddressValidator
	testCases := []struct {
		name           string
		addressType    string
		expectedResult domain.UserAddressValidator
		expectedError  error
	}{
		{
			name:           "addressType as email",
			addressType:    "email",
			expectedResult: EmailAddressValidator{},
			expectedError:  nil,
		},
		{
			name:           "addressType as phone",
			addressType:    "phone",
			expectedResult: nilAddressValidator,
			expectedError:  errors.New("no UserAddressValidotor exists for give addressType: phone"),
		},
		{
			name:           "addressType as empty string",
			addressType:    "",
			expectedResult: nilAddressValidator,
			expectedError:  errors.New("no UserAddressValidotor exists for give addressType: "),
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewUserAddressValidator(tc.addressType)
			if result != tc.expectedResult ||
				(err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) {

				t.Errorf("%s: NewUserAddressValidator(%s) => Got: (%v, %v), expected: (%v, %v)", tc.name, tc.addressType, result, err, tc.expectedResult, tc.expectedError)
			}
		})
	}
}
