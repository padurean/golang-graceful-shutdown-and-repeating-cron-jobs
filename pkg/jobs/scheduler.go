package jobs

import (
	"log"
	"sync"
	"time"
)

// StartJobs ...
func StartJobs(js []*Job, wg *sync.WaitGroup, stop <-chan struct{}) {
	stopChans := make([]chan struct{}, 0, len(js))
	for _, j := range js {
		stopChan := make(chan struct{}, 1)
		go j.Start(stopChan, wg)
		stopChans = append(stopChans, stopChan)
	}
	<-stop
	for _, stopChan := range stopChans {
		stopChan <- struct{}{}
		close(stopChan)
	}
}

// WaitWithTimeout waits the waitgroup for the specified time(out).
// Returns true if waiting timed out.
func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) (bool, time.Duration) {
	log.Printf("Waiting for jobs to finish any ongoing runs (timeout: %v)", timeout)

	start := time.Now()

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		took := time.Since(start)
		log.Printf("Jobs finished any ongoing runs in %s", took)
		return false, took
	case <-time.After(timeout):
		took := time.Since(start)
		log.Printf(
			"WARNING: Waiting for jobs to finish any ongoing runs timed out after %s", took)
		return true, took
	}
}
