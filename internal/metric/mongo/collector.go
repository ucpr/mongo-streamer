package mongo

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace is the namespace for the metrics.
	namespace = "mongo_streamer"
	// subSystem is the subSystem for the metrics.
	subSystem = "mongodb"

	lDatabase   = "database"
	lCollection = "collection"
)

var (
	// receivedTotal is the total number of change stream received.
	receivedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "change_stream_received_total",
			Help:      "Total number of change stream received",
		}, []string{lDatabase, lCollection},
	)

	// receivedBytesTotal is the total number of change stream received bytes.
	receivedBytesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "change_stream_received_bytes_total",
			Help:      "Total number of change stream received bytes",
		}, []string{lDatabase, lCollection},
	)

	// successHandleEventTotal is the total number of change stream handle event success.
	successHandleEventTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "change_stream_handle_event_success_total",
			Help:      "Total number of change stream handle event success",
		}, []string{lDatabase, lCollection},
	)

	// failedHandleEventTotal is the total number of change stream handle event failed.
	failedHandleEventTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "change_stream_handle_event_failed_total",
			Help:      "Total number of change stream handle event failed",
		}, []string{lDatabase, lCollection},
	)
)

// Collectors returns all collectors of MongoDB.
func Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		receivedTotal,
		receivedBytesTotal,
		successHandleEventTotal,
		failedHandleEventTotal,
	}
}

// ReceiveChangeStream increase the total number of change stream received.
func ReceiveChangeStream(database, collection string) {
	receivedTotal.WithLabelValues(database, collection).Inc()
}

// ReceiveChangeStreamBytes increase the total number of change stream received bytes.
func ReceiveBytes(database, collection string, size int) {
	receivedBytesTotal.WithLabelValues(database, collection).Add(float64(size))
}

// HandleChangeEventSuccess increase the total number of change stream handle event success.
func HandleChangeEventSuccess(database, collection string) {
	successHandleEventTotal.WithLabelValues(database, collection).Inc()
}

// HandleChangeEventFailed increase the total number of change stream handle event failed.
func HandleChangeEventFailed(database, collection string) {
	failedHandleEventTotal.WithLabelValues(database, collection).Inc()
}
