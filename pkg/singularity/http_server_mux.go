package singularity

import (
	"net/http"

	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
)

type healthcheck struct {
	failing bool
	reason  string
}

func (h healthcheck) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if h.failing {
		rw.WriteHeader(http.StatusServiceUnavailable)
		rw.Write([]byte(h.reason))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}

func WithHTTPServerMux(s *Singularity) cloudeventshttp.Option {
	return func(t *cloudeventshttp.Transport) error {
		s.healthcheck = &healthcheck{}

		mux := http.NewServeMux()
		mux.Handle("/healthz", s.healthcheck)

		t.Handler = mux

		return nil
	}
}
