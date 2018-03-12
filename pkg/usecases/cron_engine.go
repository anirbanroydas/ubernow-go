package usecases

import (
	// "fmt"
	"time"
)

type CronEngine interface {
	Add(time.Time, CronJob) error
}

type CronJob func()
