package usecases

import (
	// "fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type Job interface {
	DoWork() error
}

type AppEngine interface {
	AddJob(Job) error
}

// UserRequestJob implements the Job interface.
type UserRequestJob struct {
	UserRequest                   *domain.UserRequest
	TrafficInteractor             *TrafficInteractor
	CabInteractor                 *CabInteractor
	CabEngineInteractor           *CabEngineInteractor
	NotificationInteractor        *NotificationInteractor
	NotificationServiceInteractor *NotificationServiceInteractor
	CronEngine                    CronEngine
}

func (job *UserRequestJob) DoWork() error {
	var err error
	var baseTravelTime time.Time
	baseTravelTime, err = job.TrafficInteractor.GetBaseTravelTime(job.UserRequest.Source, job.UserRequest.Destination, time.Now())
	if err != nil {
		return errors.Wrap(err, "UserRequestJob's DoWork returned error while calling TrafficInteractor.GetBaseTravelTime method")
	}

	var baseEta time.Duration
	baseEta, err = job.CabInteractor.GetBaseEta(job.UserRequest.Source, job.UserRequest.Destination, time.Now(), job.UserRequest.Cab, job.UserRequest.CabType)
	if err != nil {
		return errors.Wrap(err, "UserRequestJob's DoWork returned error while calling CabInteractor.GetBaseEta method")
	}

	var tResp *TrafficResponseDTO
	tResp, err = job.TrafficInteractor.GetTrafficFinalResponse(baseTravelTime, job.UserRequest)
	if err != nil {
		return errors.Wrap(err, "UserRequestJob's DoWork returned error while calling TrafficInteractor.GetTrafficFinalResponse method")
	}

	// step 3: pass the result to cab request's job queue,(but not directly).
	// First, create cron jobs which will trigger the functions at the specific time
	// which will then add those jobs to the cab request's job queu, which will be consumed by
	// cab request specific worker which will then find the final booking time for the cab
	triggerTime := job.TrafficInteractor.GetTriggerTime(baseEta, tResp)
	err = job.CronEngine.Add(triggerTime, job.CabEngineInteractor.TrafficResponseProcessor(tResp, job.CabInteractor, job.NotificationInteractor, job.NotificationServiceInteractor))
	if err != nil {
		return errors.Wrap(err, "UserRequestJob's DoWork couldn't perform Cron.Add()")
	}
	return nil
}

func NewUserRequestJob(ur *domain.UserRequest, tsI *TrafficInteractor, cabI *CabInteractor, cabEngI *CabEngineInteractor, n *NotificationInteractor, c CronEngine) *UserRequestJob {
	job := UserRequestJob{
		UserRequest:            ur,
		TrafficInteractor:      tsI,
		CabInteractor:          cabI,
		CabEngineInteractor:    cabEngI,
		NotificationInteractor: n,
		CronEngine:             c,
	}
	return &job
}

// TrafficAppEngine implements the AppEngine interface to create a Job of type
// UserRequestJob and add that job to a JobQueue which is then consumed by
// some workers to perform all the traffic_request use cases.
type TrafficAppEngine struct {
	JobQueue chan Job
}

func (t *TrafficAppEngine) AddJob(j Job) error {
	select {
	case t.JobQueue <- j:
		return nil
	default:
		return errors.New("TrafficAppEngine Couldn't add any more jobs to JobQueue")
	}

}

func NewTrafficAppEngine(maxQueueLength int) *TrafficAppEngine {
	t := TrafficAppEngine{
		JobQueue: make(chan Job, maxQueueLength),
	}

	return &t
}

// CabRequestJob implements the Job interface
type CabRequestJob struct {
	TrafficResponse               *TrafficResponseDTO
	CabInteractor                 *CabInteractor
	NotificationInteractor        *NotificationInteractor
	NotificationServiceInteractor *NotificationServiceInteractor
}

func (job *CabRequestJob) DoWork() error {
	bResp, err := job.CabInteractor.GetBookingResponse(job.TrafficResponse)
	if err != nil {
		return errors.Wrap(err, "CabRequestJob's DoWork errored while calling GetBookingResponse")
	}

	// step 1: Use the NotificationInteractor to send the booking respone
	// to notification service via its SendToQueue method.
	// NOTE: NotificationInteractor uses a queue again to immediately send
	// the CabBookingResponse to where a worker will pick it up and do the final processing.
	// This is done to decouple the responsibility
	// and also for asynchronous behaviour.
	err = job.NotificationInteractor.SendQueue(bResp, job.NotificationServiceInteractor)
	if err != nil {
		return errors.Wrap(err, "CabRequestJob's DoWork method returned error while calling SendToQueue method of NotificationInteractor")
	}

	return nil
}

func NewCabRequestJob(t *TrafficResponseDTO, c *CabInteractor, nI *NotificationInteractor, nsI *NotificationServiceInteractor) *CabRequestJob {
	job := CabRequestJob{
		TrafficResponse:               t,
		CabInteractor:                 c,
		NotificationInteractor:        nI,
		NotificationServiceInteractor: nsI,
	}
	return &job
}

// CabAppEngine implements the AppEngine interface to create a Job of type
// CabRequestJob and add that job to a JobQueue which is then consumed by
// some workers to perform all the traffic_request use cases.
type CabAppEngine struct {
	JobQueue chan Job
}

func (c *CabAppEngine) AddJob(j Job) error {
	select {
	case c.JobQueue <- j:
		return nil
	default:
		return errors.New("CabAppEngine Couldn't add any more jobs to JobQueue")
	}
}

func NewCabAppEngine(maxQueueLength int) *CabAppEngine {
	c := CabAppEngine{
		JobQueue: make(chan Job, maxQueueLength),
	}

	return &c
}

// NotificationJob implements the Job interface
type NotificationJob struct {
	*domain.CabBookingResponse
	NotificationServiceInteractor *NotificationServiceInteractor
}

func (job *NotificationJob) DoWork() error {
	err := job.NotificationServiceInteractor.Send(job.CabBookingResponse)
	if err != nil {
		return errors.Wrap(err, "NotificationJob's DoWork method returned error while calling Send method of NotificationServiceInteractor")
	}

	return nil
}

func NewNotificationJob(c *domain.CabBookingResponse, n *NotificationServiceInteractor) *NotificationJob {
	job := NotificationJob{
		CabBookingResponse:            c,
		NotificationServiceInteractor: n,
	}
	return &job
}

// NotificationAppEngine implements the AppEngine interface to create a Job of type
// CabRequestJob and add that job to a JobQueue which is then consumed by
// some workers to perform all the traffic_request use cases.
type NotificationAppEngine struct {
	JobQueue chan Job
}

func (n *NotificationAppEngine) AddJob(j Job) error {
	select {
	case n.JobQueue <- j:
		return nil
	default:
		return errors.New("NotificationAppEngine Couldn't add any more jobs to JobQueue")
	}
}

func NewNotificationAppEngine(maxQueueLength int) *NotificationAppEngine {
	n := NotificationAppEngine{
		JobQueue: make(chan Job, maxQueueLength),
	}

	return &n
}
