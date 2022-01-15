package jobs_test

import (
	"sync"
	"testing"
	"time"

	"github.com/padurean/golang-graceful-shutdown-and-repeating-cron-jobs/pkg/jobs"
)

func TestJobs(t *testing.T) {
	js := []*jobs.Job{
		{
			Name:         "Job One",
			Run:          func() { time.Sleep(250 * time.Millisecond) },
			InitialDelay: 0,
			Interval:     750 * time.Second,
		},
		{
			Name:         "Job Two",
			Run:          func() { time.Sleep(500 * time.Millisecond) },
			InitialDelay: 500 * time.Millisecond,
			Interval:     1 * time.Second,
		},
	}

	var jobsWaitGroup sync.WaitGroup
	jobsWaitGroup.Add(len(js))

	stopJobsChan := make(chan struct{}, 1)
	defer close(stopJobsChan)

	go jobs.StartJobs(
		js,
		&jobsWaitGroup,
		stopJobsChan)

	time.Sleep(1 * time.Second)

	stopJobsChan <- struct{}{}
	jobs.WaitWithTimeout(&jobsWaitGroup, 1*time.Second)
}
