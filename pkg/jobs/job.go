package jobs

import (
	"log"
	"sync"
	"time"
)

// Job repeatedly runs the specified job (function) at the specified interval.
// An optional initial delay can be specified: if > 0, the first run will happen
// after the specified delay (usually it is smaller than the interval, but not necessarily).
type Job struct {
	Name         string        // job name
	Run          func()        // job to run
	InitialDelay time.Duration // initial delay before first run (optional)
	Interval     time.Duration // interval between runs

	startOnlyOnce sync.Once
}

// Start starts the job (scheduled runs).
// wg arg is an optional wait group for signaling that the stop signal has
// been received and any ongoing runs have completed.
func (j *Job) Start(stopChan <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf(
		"Starting job '%s' (initial delay: %s, interval: %s)",
		j.Name, j.InitialDelay, j.Interval)

	if j.Interval <= 0 {
		panic("job interval must be greater than 0")
	}

	j.startOnlyOnce.Do(func() {

		var ticker *time.Ticker

		defer func() {
			if ticker != nil {
				ticker.Stop()
			}
			if wg != nil {
				wg.Done()
			}
		}()

		interval := j.InitialDelay
		initialRun := true
		if interval == 0 {
			interval = j.Interval
			initialRun = false
		}

		ticker = time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				j.run(initialRun)
				if initialRun {
					ticker.Stop()
					ticker = time.NewTicker(j.Interval)
					initialRun = false
				}
			case <-stopChan:
				log.Printf("Job '%s' stopped at %v", j.Name, time.Now())
				return
			}
		}

	})
}

func (j *Job) run(initialRun bool) {
	ir := ""
	if initialRun {
		ir = "(initial run) "
	}
	start := time.Now()
	log.Printf("Job '%s' %srunning at %v", j.Name, ir, start)

	j.Run()

	end := time.Now()
	log.Printf(
		"Job '%s' %sfinished at %v took %s", j.Name, ir, time.Now(), end.Sub(start))
}
