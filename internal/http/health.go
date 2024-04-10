package http

import (
	"net/http"

	"github.com/ucpr/mongo-streamer/pkg/log"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response: %v", err)
	}
}
