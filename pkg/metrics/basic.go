package metrics

import (
	"context"
	"fmt"
	"net/http"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Basic struct {
	server *http.Server
}

func (m *Basic) Listen(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		log.Info().
			Str("address", m.server.Addr).
			Msg("Stopping Metrics")

		m.server.Shutdown(ctx)
	}()

	log.Info().
		Str("address", m.server.Addr).
		Msg("Starting Metrics")

	err := m.server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (m *Basic) Record(event cloudevents.Event, err error) error {
	switch err == nil {
	case true:
		eventsProcessed.
			WithLabelValues(event.Type()).
			Inc()

	case false:
		eventsError.
			WithLabelValues(err.Error()).
			Inc()
	}

	return nil
}

func NewBasic(e v1alpha2.Metrics) (Metrics, error) {
	router := http.NewServeMux()
	router.Handle(e.Path, promhttp.Handler())

	return &Basic{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", e.Port),
			Handler: router,
		},
	}, nil
}
