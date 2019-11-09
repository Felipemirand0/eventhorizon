package metrics

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics interface {
	Listen(context.Context) error
	Record(cloudevents.Event, error) error
}

var eventsProcessed = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "events_processed_total",
		Help: "The total number of processed events",
	},
	[]string{"type"},
)

var eventsError = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "events_error_total",
		Help: "The total number of failed events",
	},
	[]string{"error"},
)
