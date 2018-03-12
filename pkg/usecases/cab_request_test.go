package usecases

import (
	// "fmt"
	// "reflect"
	"testing"
	"time"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type MockCabService struct{}

func (c *MockCabService) EtaNow(cr *domain.CabRequest) (time.Duration, error) {
	return time.Duration(7 * time.Minute), nil
}

type MockBookingTimeFinder struct{}

func (s *MockBookingTimeFinder) FindBest(tr *TrafficResponseDTO) (time.Time, error) {
	return time.Now().Add(5 * time.Hour), nil
}

func testCabInteractor(t *testing.T) *CabInteractor {
	t.Helper()

	cs := &MockCabService{}
	s := &MockBookingTimeFinder{}
	return NewCabInteractor(cs, s)
}

func testCabEngineInteractor(t *testing.T) *CabEngineInteractor {
	t.Helper()

	a := &MockAppEngine{}
	l := &MockLogger{}
	return NewCabEngineInteractor(a, l)
}
