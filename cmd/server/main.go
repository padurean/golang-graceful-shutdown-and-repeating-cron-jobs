package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/padurean/golang-graceful-shutdown-and-repeating-cron-jobs/internal/http/server"
	"github.com/padurean/golang-graceful-shutdown-and-repeating-cron-jobs/pkg/jobs"
)

func main() {
	errChan, err := run()
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errChan; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run() (<-chan error, error) {
	// load config from env
	port := 8080

	// do any needed initializations here - e.g. connect to db, etc.
	// and extend the httpserver.New function to accept extra args for them

	// initialize server
	srv := server.New(port)

	// start cron jobs
	stopJobsChan := make(chan struct{}, 1)
	var jobsWaitGroup sync.WaitGroup
	startJobs(&jobsWaitGroup, stopJobsChan)

	// create main error channel
	errChan := make(chan error, 1)

	// configure listening to OS signals for shutting down gracefully
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// receive shutdown signal (e.g. Ctrl+C)
		<-ctx.Done()
		log.Println("Shutdown signal received")

		// stop jobs and wait for any ongoing runs to finish
		stopJobsChan <- struct{}{}
		jobs.WaitWithTimeout(&jobsWaitGroup, 30*time.Second)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// close resources, cancel contexts and close the main error channel
		defer func() {
			// do any other cleanup here - e.g. close the database connection, etc.
			stop()
			cancel()
			close(stopJobsChan)
			close(errChan)
		}()

		// shutdown the HTTP server
		srv.SetKeepAlivesEnabled(false)
		log.Println("Shutting down the HTTP server")
		if err := srv.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		log.Println("Shutdown completed")
	}()

	go func() {
		// start the HTTP server
		log.Println("Starting HTTP server on port", port)
		// ListenAndServe always returns a non-nil error.
		// After Shutdown or Close, the returned error is ErrServerClosed.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	return errChan, nil
}

func startJobs(wg *sync.WaitGroup, stop <-chan struct{}) {
	operaJobs := []*jobs.Job{
		{
			Name: "Job One",
			Run: func() {
				log.Println("Job One running ...")
				time.Sleep(3 * time.Second)
				log.Println("Job One done")
			},
			InitialDelay: 5 * time.Second,
			Interval:     10 * time.Second,
		},
		{
			Name: "Job Two",
			Run: func() {
				log.Println("Job Two running ...")
				time.Sleep(5 * time.Second)
				log.Println("Job Two done")
			},
			InitialDelay: 10 * time.Second,
			Interval:     15 * time.Second,
		},
	}
	wg.Add(len(operaJobs))
	go jobs.StartJobs(operaJobs, wg, stop)
}
