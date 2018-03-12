package usecases

import (
	// "fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

// UserRequestDTO is DTO which takes in a UserRequest objec
type UserRequestDTO struct {
	name             string
	source           domain.Location
	destination      domain.Location
	reachingTime     time.Time
	cab              string
	cabType          string
	notificationAddr domain.UserAddress
}

type UserInteractor struct {
	UserRepository                domain.UserRepository
	RequestRepository             domain.RequestRepository
	AppEngine                     AppEngine
	// CronEngine                    CronEngine
	// TrafficInteractor             *TrafficInteractor
	// CabInteractor                 *CabInteractor
	// CabEngineInteractor           *CabEngineInteractor
	// NotificationInteractor        *NotificationInteractor
	// NotificationServiceInteractor *NotificationServiceInteractor
}

// CreateUserRequest use_case takes a UserRequestDTO object as input and creates a domain level
// UserTequest object and sends it to an AppEngine which process it from there on, asynchronously.
// It returns an error if there is a problem in any of the above processes.
//
// To do it it firest cretaes a domain level REquest object and stores it in the RequestRepository,
// then creates a domain level User Object and stores it in the
// UserRepository and then it creates the domain level UserRequest object and sends it to the AppEngine.
func (ur *UserInteractor) CreateUserRequest(ucReq UserRequestDTO) error {
	var err error

	// step 1:  create  new domoan.Request and save in domain.RequestRepository
	var r *domain.Request
	r, err = ur.createAndSaveRequest(ucReq)
	if err != nil {
		return errors.Wrap(err, "CreateUserRequest couldn't create and save domain.Request")
	}

	// step 2: create and save domain.User in domain.UserRepository
	var u *domain.User
	u, err = ur.createAndSaveUser(ucReq.name)
	if err != nil {
		return errors.Wrap(err, "CreateUserRequest couldn't create and save domain.User")
	}

	// step 3: create new domain.UserRequest Object
	userRequest := domain.NewUserRequest(u, r)

	// stpe 4: send the new domain.UserRequest to the app engine to process and return
	err = ur.sendQueue(userRequest)
	if err != nil {
		return errors.Wrap(err, "CreateUserRequest could't send userRequest to AppEngine for processing")
	}

	return nil
}

// createAndSaveRequest is a method of UserInteractor struct which takes in a UserRequestDTO object as input
// an creates a domain.Request object and stores it in the domain.UserRepository.
func (ur *UserInteractor) createAndSaveRequest(ucReq UserRequestDTO) (*domain.Request, error) {
	// step 1: create a new address validator by using the user's given address type
	var r *domain.Request
	uav, err := NewUserAddressValidator(ucReq.notificationAddr.AddrType)
	if err != nil {
		return r, errors.Wrap(err, "createAndSaveRequest can't creat new UserAddressValidator")
	}
	// step 2: create new domain.Request
	r, err = domain.NewRequest(ucReq.source, ucReq.destination, ucReq.reachingTime, ucReq.cab, ucReq.cabType, ucReq.notificationAddr, uav)
	if err != nil {
		return r, errors.Wrap(err, "createAndSaveRequest can't creat New domain.request Object")
	}

	// step 3: save domain.Request object to domain.RequestRepository
	_, err = ur.RequestRepository.Store(r)
	if err != nil {
		return r, errors.Wrap(err, "createAndSaveRequest couldn't store request to RequestRepository")
	}

	return r, nil
}

// createAndSaveUser is a method of UserInteractor which takes a name of type string
// as input and returns a domain.User ojbect along with an error.
func (ur *UserInteractor) createAndSaveUser(name string) (*domain.User, error) {
	var u *domain.User
	// step 1: create new domain.User object
	u = domain.NewUser(name)

	// step 2: save domain.User object in domain.UserRepository object.
	_, err := ur.UserRepository.Store(u)
	if err != nil {
		return u, errors.Wrap(err, "createAndSaveUser couldn't store user to UserRepository")
	}

	return u, nil
}

// sendQueue is a method on UserInteractor which takes a pointer to domain.UserRequest
// as input and returns an error. This method create a new UserRequestJob which is sent to
// a JobQueue which is later processed by some worker. The UserRequestJob has information
// about what to do and how to do as there are injected into the UserRequestJob object.
func (ur *UserInteractor) sendQueue(userRequest *domain.UserRequest) error {
	// step 1: create a new UserRequestJob which is Job interface
	job := NewUserRequestJob(userRequest)

	// step 2: add the new job to AppEngine Queue
	err := ur.AppEngine.AddJob(job)
	if err != nil {
		return errors.Wrap(err, "sendAppEngine couldn't add the new  AppEngine Job for the userRequest")
	}

	return nil
}

// NewUserInteractor is consturctor
func NewUserInteractor(uRepo domain.UserRepository, reqRepo domain.RequestRepository, c CronEngine, a AppEngine, trI *TrafficInteractor, cabI *CabInteractor, cabEngI *CabEngineInteractor, nI *NotificationInteractor, nsI *NotificationServiceInteractor) *UserInteractor {
	u := UserInteractor{
		UserRepository:                uRepo,
		RequestRepository:             reqRepo,	
		AppEngine:                     a,
		// CronEngine:                    c,
		// TrafficInteractor:             trI,
		// CabInteractor:                 cabI,
		// CabEngineInteractor:           cabEngI,
		// NotificationInteractor:        nI,
		// NotificationServiceInteractor: nsI,
	}
	return &u
}



 ur.TrafficInteractor, ur.CabInteractor, ur.CabEngineInteractor, ur.NotificationInteractor, ur.CronEngine