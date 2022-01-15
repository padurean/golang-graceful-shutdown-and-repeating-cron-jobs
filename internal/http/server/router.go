package server

import (
	"net/http"

	"github.com/padurean/golang-graceful-shutdown-and-repeating-cron-jobs/internal/http/handlers"
)

func router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleGetHello)
	mux.HandleFunc("/ping", handlers.HandlePostPing)

	return mux
}
