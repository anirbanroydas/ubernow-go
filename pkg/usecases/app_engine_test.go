package usecases

import (
	// "fmt"
	// "reflect"
	// "testing"
	"time"

	"github.com/pkg/errors"
	// "github.com/anirbanroydas/ubernow-go/pkg/domain"
)

type MockCronEngine struct{}

func (c *MockCronEngine) Add(triggerTime time.Time, cJob CronJob) error {
	return nil
}

type MockAppEngine struct{}

func (a *MockAppEngine) AddJob(j Job) error {
	return nil
}

type MockBadAppEngine struct{}

func (a *MockBadAppEngine) AddJob(j Job) error {
	return errors.New("couldn't add job to JobQueue")
}
