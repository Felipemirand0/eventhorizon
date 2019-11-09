package eventhorizon

import (
	"net/http"

	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
)

type healthcheck struct {
	failing bool
	reason  error
}

func (h healthcheck) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if h.failing {
		rw.WriteHeader(http.StatusServiceUnavailable)
		rw.Write([]byte(h.reason.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}

func WithHTTPServerMux(h *healthcheck) cloudeventshttp.Option {
	return func(t *cloudeventshttp.Transport) error {
		mux := http.NewServeMux()
		mux.Handle("/healthz", h)

		t.Handler = mux

		return nil
	}
}
