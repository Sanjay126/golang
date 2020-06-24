package handler

import (
	"bitbucket.org/dreamplug-backend/benchmarks/stats"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const serviceName = "demo-app"

// WithStats wraps handlers with stats reporting. It tracks metrics such
// as the number of requests per endpoint, the latency, etc.
func WithStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := getStatsTags(r)
		stats.IncCounter("handler.received", tags, 1)
		h(w, r)
	}
}
var host, _= os.Hostname()
func getStatsTags(r *http.Request) map[string]string {
	statsTags := map[string]string{
		"endpoint": filepath.Base(r.URL.Path),
	}

	if idx := strings.IndexByte(host, '.'); idx > 0 {
		host = host[:idx]
	}
	statsTags["host"] = host
	return statsTags
}

