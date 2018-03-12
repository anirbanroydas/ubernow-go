package usecases

import (
	// "fmt"
	// "reflect"
	"testing"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

func testNotificationInteractor(t *testing.T) *NotificationInteractor {
	t.Helper()

	a := &MockAppEngine{}
	return NewNotificationInteractor(a)
}

type MockNotificationService struct{}

func (n *MockNotificationService) Send(c *domain.CabBookingResponse) error {
	return nil
}

func testNotificationServiceInteractor(t *testing.T) *NotificationServiceInteractor {
	t.Helper()

	ns := &MockNotificationService{}
	return NewNotificationServiceInteractor(ns)
}
