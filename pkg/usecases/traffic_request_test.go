package usecases

import (
	// "fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type MockTrafficService struct{}

func (t *MockTrafficService) TravelTime(tr *domain.TrafficRequest) (*domain.TrafficResponse, error) {
	tResp := domain.TrafficResponse{
		TrafficRequest: tr,
		TravelTime:     time.Duration(45 * time.Minute),
		BestCase:       time.Now().Add(4 * time.Hour),
		WorstCase:      time.Now().Add(210 * time.Minute),
	}

	return &tResp, nil
}

func testTrafficInteractor(t *testing.T) *TrafficInteractor {
	t.Helper()

	ts := &MockTrafficService{}
	return NewTrafficInteractor(ts)
}
