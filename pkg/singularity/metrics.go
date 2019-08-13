package singularity

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricEventsProcessed = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "events_processed_total",
		Help: "The total number of processed events",
	},
	[]string{"type"},
)
