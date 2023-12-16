package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ucpr/mongo-streamer/internal/metric/mongo"
)

// Register register prometheus metrics to http.ServeMux
func Register(mux *http.ServeMux) {
	reg := prometheus.NewRegistry()

	// register metrics
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
	reg.MustRegister(
		mongo.Collectors()...,
	)

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
}
