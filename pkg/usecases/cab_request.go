package usecases

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type CabInteractor struct {
	CabService domain.CabService
	Strategy   BestBookingTimeFinder
}

type CabEngineInteractor struct {
	AppEngine AppEngine
	Logger    domain.Logger
}

type BestBookingTimeFinder interface {
	FindBest(*TrafficResponseDTO) (time.Time, error)
}

type HeuristicBestTimeStrategy struct {
}

func (h *HeuristicBestTimeStrategy) FindBest(tr *TrafficResponseDTO) (time.Time, error) {
	var bestTime time.Time
	// only considering bestCase response
	// TODO: Implement this algo
	// step 0: pool cab service for current eta
	// step 1: check if time already passed with the eta and the
	// last travel time from traffic response
	// step 2: find difference with all the possible traffic
	// response's worst and best case times
	// step 3: calculate avg eta and continue to poll cab service until
	// you reach a eta which is closes to any best/worst case travel time
	// from traffic responses by a threshold amount
	// step 4: craete a cabBookingResponse object and return it

	// for now return current time
	bestTime = time.Now() // Fix This
	return bestTime, nil
}

func (c *CabInteractor) GetBaseEta(source, destination domain.Location, bookingTime time.Time, cab, cabType string) (time.Duration, error) {
	var baseEta time.Duration
	// step 0: create new cab request
	cabReq := domain.NewCabRequest(source, destination, bookingTime, cab, cabType)
	// step 1: poll traffic service
	eta, err := c.CabService.EtaNow(cabReq)
	if err != nil {
		return baseEta, errors.Wrap(err, "GetBaseEta failed in fetching EtaNow from CabService")
	}
	baseEta = eta

	return baseEta, nil
}

func (c *CabInteractor) GetBookingResponse(tr *TrafficResponseDTO) (*domain.CabBookingResponse, error) {
	var cResp *domain.CabBookingResponse
	// Use the strategy which is associated with this CabServiceIndicator to find the BestTime
	// possible
	bestBookingTime, err := c.Strategy.FindBest(tr)
	if err != nil {
		return cResp, errors.Wrap(err, "CabInteractor's GetBookingResponse returned error while calling its Strategy's FindBest method")
	}

	// create the CabBookingResponse object
	cResp = domain.NewCabBookingResponse(tr.UserRequest, bestBookingTime)
	return cResp, nil
}

func (c *CabEngineInteractor) TrafficResponseProcessor(tr *TrafficResponseDTO, cs *CabInteractor, nI *NotificationInteractor, nsI *NotificationServiceInteractor) CronJob {
	return func() {
		err := c.sendQueue(tr, cs, nI, nsI)
		if err != nil {
			c.Logger.LogError(fmt.Sprintf("CabEngineInteractor.sendToQueue Error:: %v", err))
		}
	}
}

func (c *CabEngineInteractor) sendQueue(tr *TrafficResponseDTO, cs *CabInteractor, nI *NotificationInteractor, nsI *NotificationServiceInteractor) error {
	// step 1: create a new app engine job
	job := NewCabRequestJob(tr, cs, nI, nsI)

	// step 2: add the new job to AppEngine Queue
	err := c.AppEngine.AddJob(job)
	if err != nil {
		return errors.Wrap(err, "sendToQueue couldn't add the new  AppEngine Job for the userRequest")
	}

	return nil
}

func NewCabEngineInteractor(a AppEngine, l domain.Logger) *CabEngineInteractor {
	c := CabEngineInteractor{
		AppEngine: a,
		Logger:    l,
	}
	return &c
}

func NewCabInteractor(cs domain.CabService, s BestBookingTimeFinder) *CabInteractor {
	c := CabInteractor{
		CabService: cs,
		Strategy:   s,
	}
	return &c
}

// TODO: imnplement these private methods to be used in the public method GetBookingResponse
func (c *CabInteractor) pollCabService() {

}

func (c *CabInteractor) sendRequestCabService() {

}

func (c *CabInteractor) makeDecision() {

}

func (c *CabInteractor) goodHeuristic() {

}

func (c *CabInteractor) sendJobToNotificationService() {

}
