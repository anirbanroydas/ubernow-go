package usecases

import (
	// "fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type TrafficResponseDTO struct {
	*domain.UserRequest
	TravelTime []time.Duration
	BestCase   []time.Time
	WorstCase  []time.Time
}

type TrafficInteractor struct {
	TrafficService domain.TrafficService
}

func (tr *TrafficInteractor) GetTrafficFinalResponse(baseTravelTime time.Duration, ur *domain.UserRequest) (*TrafficResponseDTO, error) {
	var tResp *TrafficResponseDTO
	// TODO: implement the algorithm and find the final set of
	// best and worst case times to start the journey to reach
	// destination and send the result to another job queue
	// from where the cab requet use case will consume the jobs
	// and them perform their processing to find the right time to
	// book a cab

	// step o: create new traffic request
	// step 1: poll traffic service for a traffic response
	// step 2: parse response and check for result
	// step 3: continue polling with new trafficRequest
	// until some threshold  number of time or stop on
	// reaching a condition

	// step 4: create thet traffic response dto and return
	return tResp, nil

}

func (tr *TrafficInteractor) GetBaseTravelTime(source, destination domain.Location, t time.Time) (time.Duration, error) {
	var baseTravelTime time.Duration
	// step 0: create new traffic request
	treq := domain.NewTrafficRequest(source, destination, t)
	// step 1: poll traffic service
	tresp, err := tr.TrafficService.TravelTime(treq)
	if err != nil {
		return baseTravelTime, errors.Wrap(err, "GetBaseTravelTime failed in fetching TravelTime from TrafficService")
	}
	baseTravelTime = tresp.TravelTime

	return baseTravelTime, nil
}

func (tr *TrafficInteractor) GetTriggerTime(baseEta time.Duration, tResp *TrafficResponseDTO) time.Time {
	// step 0: sort the tResp.WorstCase and find the smalled time
	// step 2: return the the difference for step0 and baseEta

	// TODO: // now just return the current time
	return time.Now()
}

func NewTrafficInteractor(ts domain.TrafficService) *TrafficInteractor {
	t := TrafficInteractor{
		TrafficService: ts,
	}
	return &t
}
