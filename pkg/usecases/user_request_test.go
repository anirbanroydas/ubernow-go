package usecases

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

// MockUserRepo implememts the domain.RequestRepository interface which alwasy returns nil error
type MockUserRepo struct{}

func (ur *MockUserRepo) Store(u *domain.User) (uint64, error) {
	var MockID uint64
	MockID = 123
	return MockID, nil
}

func (ur *MockUserRepo) FindByID(userID uint64) (*domain.User, error) {
	MockUser := domain.NewUser("roy")
	return MockUser, nil
}

// MockBadUserRepo implememts the domain.RequestRepository interface which alwasy returns some non nil error
type MockBadUserRepo struct{}

func (ur *MockBadUserRepo) Store(u *domain.User) (uint64, error) {
	var MockID uint64
	MockID = 123
	return MockID, errors.New("couldn't store user")
}

func (ur *MockBadUserRepo) FindByID(userID uint64) (*domain.User, error) {
	var MockUser *domain.User
	return MockUser, errors.New("couldn't find user")
}

// MockRequestRepo implememts the domain.RequestRepository interface which alwasy returns nil error
type MockRequestRepo struct{}

func (rp *MockRequestRepo) Store(r *domain.Request) (uint64, error) {
	var MockID uint64
	MockID = 456
	return MockID, nil
}

func (rp *MockRequestRepo) FindByID(reqID uint64) (*domain.Request, error) {
	MockUser := &domain.Request{}
	return MockUser, nil
}

// MockBadRequestRepo implememts the domain.RequestRepository interface which alwasy returns errors
type MockBadRequestRepo struct{}

func (rp *MockBadRequestRepo) Store(r *domain.Request) (uint64, error) {
	var MockID uint64
	MockID = 456
	return MockID, errors.New("couldn't store request to repo")
}

func (rp *MockBadRequestRepo) FindByID(reqID uint64) (*domain.Request, error) {
	var MockUser *domain.Request
	return MockUser, errors.New("couldn't store request to repo")
}

// MockLogger implements the domain.Logger interface
type MockLogger struct{}

func (l *MockLogger) LogError(m string) {
	fmt.Println(m)
}

func (l *MockLogger) LogInfo(m string) {
	fmt.Println(m)
}

func testUserInteractor(t *testing.T) *UserInteractor {
	t.Helper()

	uRepo := &MockUserRepo{}
	reqRepo := &MockRequestRepo{}
	c := &MockCronEngine{}
	a := &MockAppEngine{}
	trI := testTrafficInteractor(t)
	cabI := testCabInteractor(t)
	cabEngI := testCabEngineInteractor(t)
	nI := testNotificationInteractor(t)
	nsI := testNotificationServiceInteractor(t)

	return NewUserInteractor(uRepo, reqRepo, c, a, trI, cabI, cabEngI, nI, nsI)
}

func TestCreateAndSaveUser(t *testing.T) {
	uRoy := domain.NewUser("roy")
	uEmtpy := domain.NewUser("")
	interactor := testUserInteractor(t)

	testCases := []struct {
		name          string
		userName      string
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:          "valid name and no error from user repository",
			userName:      "roy",
			expectedUser:  uRoy,
			expectedError: nil,
		},
		{
			name:          "empty name and no error from user repository",
			userName:      "",
			expectedUser:  uEmtpy,
			expectedError: nil,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			user, err := interactor.createAndSaveUser(tc.userName)
			if !reflect.DeepEqual(user, tc.expectedUser) ||
				(err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) {

				t.Errorf("%s: createAndSaveUser(%s) => got: (%v, %v) expected: (%v, %v)", tc.name, tc.userName, user, err, tc.expectedUser, tc.expectedError)
			}
		})
	}

	// Mock userepo which returns error
	interactor.UserRepository = &MockBadUserRepo{}

	testCases = []struct {
		name          string
		userName      string
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:          "valid name and some error from user repository",
			userName:      "roy",
			expectedUser:  uRoy,
			expectedError: errors.New("Some User repo store error"),
		},
		{
			name:          "empty name and some error from user repository",
			userName:      "",
			expectedUser:  uEmtpy,
			expectedError: errors.New("Some User repo store error"),
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			user, err := interactor.createAndSaveUser(tc.userName)
			if !reflect.DeepEqual(user, tc.expectedUser) ||
				(err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) {

				t.Errorf("%s: createAndSaveUser(%s) => got: (%v, %v) expected: (%v, %v)", tc.name, tc.userName, user, err, tc.expectedUser, tc.expectedError)
			}
		})
	}
}

func TestCreateAndSaveRequest(t *testing.T) {
	// base good UserRequestDTO to be sent as input to createAndSaveRequest method
	uReqDTO := UserRequestDTO{
		name: "roy",
		source: domain.Location{
			Latitude:  "77.134134",
			Longitude: "45.1341324",
		},
		destination: domain.Location{
			Latitude:  "77.234134",
			Longitude: "45.5641324",
		},
		reachingTime: time.Now().Add(5 * time.Hour),
		cab:          "uber",
		cabType:      "uberGo",
		notificationAddr: domain.UserAddress{
			AddrType: "email",
			Value:    "anirba.nick@gmail.com",
		},
	}
	// initialze all the invalid DTO's to test for different cases
	// invalid addressType
	uReqDTO1 := uReqDTO
	uReqDTO1.notificationAddr.AddrType = "web"

	// DTO with invalid source
	uReqDTO2 := uReqDTO
	uReqDTO2.source.Latitude = ""

	// DTO with invalid destination
	uReqDTO3 := uReqDTO
	uReqDTO3.destination.Longitude = ""

	// DTO with invalid reachingTime
	uReqDTO4 := uReqDTO
	uReqDTO4.reachingTime = time.Now().Add(3 * time.Minute)

	// DTO with invalid cab
	uReqDTO5 := uReqDTO
	uReqDTO5.cab = "ola"

	// DTO with invalid cabType
	uReqDTO6 := uReqDTO
	uReqDTO6.cabType = "micro"

	// DTO with invalid userAddress value
	uReqDTO7 := uReqDTO
	uReqDTO7.notificationAddr.Value = "anirban.nick@gmail"

	var emptyRequest *domain.Request
	// create a mock UserAddressValidator
	uav := MockAddressValidator{}
	// create a valid domain.Request
	r, _ := domain.NewRequest(uReqDTO.source, uReqDTO.destination, uReqDTO.reachingTime, uReqDTO.cab, uReqDTO.cabType, uReqDTO.notificationAddr, uav)
	someError := errors.New("some error")

	// initialzie the test UserRequestInteractor
	interactor := testUserInteractor(t)

	testCases := []struct {
		name            string
		uReqDTO         UserRequestDTO
		expectedRequest *domain.Request
		expectedError   error
	}{
		{
			name:            "valid uReqDTO with error in userAddressValidator",
			uReqDTO:         uReqDTO1,
			expectedRequest: emptyRequest,
			expectedError:   errors.New("no UserAddressValidotor exists for give addressType: web"),
		},
		{
			name:            "valid uReqDTO with invalid source",
			uReqDTO:         uReqDTO2,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with invalid destination",
			uReqDTO:         uReqDTO3,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with invalid reachingTime",
			uReqDTO:         uReqDTO4,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with invalid cab",
			uReqDTO:         uReqDTO5,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with invalid cabType",
			uReqDTO:         uReqDTO6,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with invalid notificationAddr",
			uReqDTO:         uReqDTO7,
			expectedRequest: emptyRequest,
			expectedError:   someError,
		},
		{
			name:            "valid uReqDTO with valid parameters and no error in RequestRepo",
			uReqDTO:         uReqDTO,
			expectedRequest: r,
			expectedError:   nil,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := interactor.createAndSaveRequest(tc.uReqDTO)
			if !reflect.DeepEqual(req, tc.expectedRequest) ||
				(err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) {

				t.Errorf("%s: createAndSaveRequest(%v) => got: (%v, %v) expected: (%v, %v)", tc.name, tc.uReqDTO, req, err, tc.expectedRequest, tc.expectedError)
			}
		})
	}

	// update interactor with a mock domain.RequestRepositoy which return error
	interactor.RequestRepository = &MockBadRequestRepo{}

	testCases = []struct {
		name            string
		uReqDTO         UserRequestDTO
		expectedRequest *domain.Request
		expectedError   error
	}{
		{
			name:            "valid uReqDTO with valid parameters but error in RequestRepo",
			uReqDTO:         uReqDTO,
			expectedRequest: r,
			expectedError:   someError,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := interactor.createAndSaveRequest(tc.uReqDTO)
			if !reflect.DeepEqual(req, tc.expectedRequest) ||
				(err != nil && tc.expectedError == nil) ||
				(err == nil && tc.expectedError != nil) {

				t.Errorf("%s: createAndSaveRequest(%v) => got: (%v, %v) expected: (%v, %v)", tc.name, tc.uReqDTO, req, err, tc.expectedRequest, tc.expectedError)
			}
		})
	}

}

func TestSendQueue(t *testing.T) {
	interactor := testUserInteractor(t)
	// create a valid domain.User
	u := domain.NewUser("roy")

	// base good UserRequestDTO to be sent as input to createAndSaveRequest method
	uReqDTO := UserRequestDTO{
		name: "roy",
		source: domain.Location{
			Latitude:  "77.134134",
			Longitude: "45.1341324",
		},
		destination: domain.Location{
			Latitude:  "77.234134",
			Longitude: "45.5641324",
		},
		reachingTime: time.Now().Add(5 * time.Hour),
		cab:          "uber",
		cabType:      "uberGo",
		notificationAddr: domain.UserAddress{
			AddrType: "email",
			Value:    "anirba.nick@gmail.com",
		},
	}
	// create a mock UserAddressValidator
	uav := MockAddressValidator{}
	// create a valid domain.Request
	r, _ := domain.NewRequest(uReqDTO.source, uReqDTO.destination, uReqDTO.reachingTime, uReqDTO.cab, uReqDTO.cabType, uReqDTO.notificationAddr, uav)
	uReq := domain.NewUserRequest(u, r)

	testCases := []struct {
		name          string
		uReq          *domain.UserRequest
		expectedError error
	}{
		{
			name:          "valid domain.UserRequest with no error from AppEngine's AddJob method",
			uReq:          uReq,
			expectedError: nil,
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			err := interactor.sendQueue(tc.uReq)
			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) {
				t.Errorf("%s: sendQueue(%v) => got: (%v) expected: (%v)", tc.name, tc.uReq, err, tc.expectedError)
			}
		})
	}

	// mocking error in AppEngine
	interactor.AppEngine = &MockBadAppEngine{}

	testCases = []struct {
		name          string
		uReq          *domain.UserRequest
		expectedError error
	}{
		{
			name:          "valid domain.UserRequest with error in AppEngine's AddJob method",
			uReq:          uReq,
			expectedError: errors.New("some eror"),
		},
	}

	for i, _ := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			err := interactor.sendQueue(tc.uReq)
			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) {
				t.Errorf("%s: sendQueue(%v) => got: (%v) expected: (%v)", tc.name, tc.uReq, err, tc.expectedError)
			}
		})
	}

}
