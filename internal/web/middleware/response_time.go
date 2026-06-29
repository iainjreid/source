package middleware

import (
	"fmt"
	"net/http"

	"github.com/iainjreid/source/internal/utils"
)

type reponseTimeReporter struct {
	http.ResponseWriter
	utils.Timeable
	wrote bool
}

func (w *reponseTimeReporter) WriteHeader(status int) {
	if !w.wrote {
		w.StopClock()
		w.Header().Set("X-Response-Time", fmt.Sprintf("%.1fms", w.TimeElapsed))
		w.wrote = true
	}

	w.ResponseWriter.WriteHeader(status)
}

func (w *reponseTimeReporter) Write(p []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(p)
}

// Timing adds the time taken to service a request as the X-Response-Time
// response header.
func ResponseTimeReporter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := &reponseTimeReporter{
			ResponseWriter: w,
		}

		tw.StartClock()
		next.ServeHTTP(tw, r)
	})
}
