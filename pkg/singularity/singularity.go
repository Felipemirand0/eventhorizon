package singularity

import (
	"context"
	"fmt"
	"sync"

	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/handler"
	"acesso.io/eventhorizon/pkg/validator"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Singularity struct {
	name       string
	mutex      *sync.RWMutex
	handlers   map[string]handler.Handler
	validators map[string]validator.Validator
	validation bool
	metrics    *metricServer
}

func (s *Singularity) handle(ctx context.Context, event cloudevents.Event) error {
	var (
		err  error
		errs []error
		lvl  zerolog.Level
	)

	defer func() {
		log.WithLevel(lvl).
			Err(err).
			Errs("errors", errs).
			Str("id", event.ID()).
			Str("type", event.Type()).
			Str("specversion", event.SpecVersion()).
			Str("source", event.Source()).
			Str("subject", event.Subject()).
			Time("time", event.Time()).
			Msg("Event received")

		if nil != s.metrics {
			go s.metrics.record(event, err)
		}
	}()

	lvl = zerolog.DebugLevel

	for key, handler := range s.handlers {
		if s.validation {
			validator := s.findHandlerValidator(key)
			if nil != validator {
				err := validator.Validate(event)
				if nil != err {
					errs = append(errs, fmt.Errorf(`failing validator "%s" for handler "%s": %s`, validator.Name(), key, err.Error()))
					continue
				}
			}
		}

		err := handler.Handle(ctx, event)
		if nil != err {
			errs = append(errs, fmt.Errorf(`failing handler "%s": %s`, key, err.Error()))
		}
	}

	if 0 == len(s.handlers) {
		err = ErrNoHandlerRegistered
		lvl = zerolog.WarnLevel

		return err
	}

	if len(errs) > 0 {
		err = ErrHandlerOneOrMoreFail
		lvl = zerolog.ErrorLevel

		return err
	}

	return nil
}

func (s *Singularity) findHandlerValidator(handlerKey string) validator.Validator {
	if nil == s.validators[handlerKey] {
		return nil
	}

	return s.validators[handlerKey]
}

func (s *Singularity) Receiver() func(context.Context, cloudevents.Event) error {
	return s.handle
}

func (s *Singularity) SetOption(opt Option) {
	opt(s)
}

func (s *Singularity) Metrics() *metricServer {
	return s.metrics
}

func New(name string, opts ...Option) (*Singularity, error) {
	s := &Singularity{
		name:       name,
		mutex:      &sync.RWMutex{},
		handlers:   map[string]handler.Handler{},
		validators: map[string]validator.Validator{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
