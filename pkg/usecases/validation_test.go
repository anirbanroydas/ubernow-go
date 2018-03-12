package usecases

import (
	// "fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

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
			t.Logf("err: %+v\n", err)
			t.Logf("tc.expectedError: %+v\n", tc.expectedError)
			if result != tc.expectedResult || !reflect.DeepEqual(err, tc.expectedError) {
				t.Errorf("%s: NewUserAddressValidator(%s) => Got: (%v, %v), expected: (%v, %v)", tc.name, tc.addressType, result, err, tc.expectedResult, tc.expectedError)
			}
		})
	}
}
