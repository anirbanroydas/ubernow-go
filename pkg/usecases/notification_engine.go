package usecases

import (
	// "fmt"

	"github.com/pkg/errors"

	"github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type NotificationInteractor struct {
	AppEngine AppEngine
}

type NotificationServiceInteractor struct {
	NotificationService domain.NotificationService
}

func (n *NotificationInteractor) SendQueue(cbResp *domain.CabBookingResponse, nsI *NotificationServiceInteractor) error {
	// step 1: create a new app engine job
	job := NewNotificationJob(cbResp, nsI)

	// step 2: add the new job to AppEngine Queue
	err := n.AppEngine.AddJob(job)
	if err != nil {
		return errors.Wrap(err, "NotificationInteractor's SendToQueue couldn't add the notificationJob via the AppEngine.AddJob method")
	}

	return nil
}

func (n *NotificationServiceInteractor) Send(c *domain.CabBookingResponse) error {
	return n.NotificationService.Send(c)
}

func NewNotificationInteractor(a AppEngine) *NotificationInteractor {
	n := NotificationInteractor{
		AppEngine: a,
	}
	return &n
}

func NewNotificationServiceInteractor(ns domain.NotificationService) *NotificationServiceInteractor {
	n := NotificationServiceInteractor{
		NotificationService: ns,
	}
	return &n
}
