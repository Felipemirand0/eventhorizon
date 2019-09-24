package singularity

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

func (s *metricServer) record(event cloudevents.Event, err error) {
	switch err == nil {
	case true:
		metricEventsProcessed.
			WithLabelValues(event.Type()).
			Inc()

	case false:
		metricEventsProcessed.
			WithLabelValues(err.Error()).
			Inc()
	}
}
