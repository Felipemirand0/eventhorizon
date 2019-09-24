package singularity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricServer struct {
	server *http.Server
}

func (s *metricServer) Listen(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.server.Shutdown(ctx)
	}()

	err := s.server.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func newMetricsServer(path string, port int) *metricServer {
	router := http.NewServeMux()
	router.Handle(path, promhttp.Handler())

	return &metricServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}
}
