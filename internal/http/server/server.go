package server

import (
	"fmt"
	"net/http"
	"time"
)

// New ...
func New(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router(),

		// set timeouts to avoid Slowloris attacks.
		ReadHeaderTimeout: time.Second * 20,
		WriteTimeout:      time.Second * 60,
		ReadTimeout:       time.Second * 60,
		IdleTimeout:       time.Second * 120,
	}
}
