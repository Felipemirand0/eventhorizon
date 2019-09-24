package singularity

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/handler"
	"acesso.io/eventhorizon/pkg/validator"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Singularity struct {
	name        string
	mutex       *sync.RWMutex
	handlers    map[string]handler.Handler
	validators  map[string]validator.Validator
	validation  bool
	metrics     *metricServer
	healthcheck *healthcheck

	backlog      int
	maxRetry     int
	retryWait    int
	maxRetryWait int

	pending chan message

	context context.Context
	cancel  context.CancelFunc
}

type message struct {
	context context.Context
	event   cloudevents.Event
}

func (s *Singularity) queue(ctx context.Context, e cloudevents.Event) error {
	s.pending <- message{
		context: ctx,
		event:   e,
	}

	return nil
}

func (s *Singularity) process() {
InfiniteLoop:
	for {
		select {
		case msg := <-s.pending:
			errs := []error{}

			handlers := s.handlers

			var wg sync.WaitGroup
			wg.Add(len(handlers))

			// do something
			for hk, h := range handlers {
				go func(hk string, h handler.Handler) {
					var err error

					for i := 0; i < s.maxRetry; i++ {
						waitTime := s.retryWait * rateinc(1.5, float64(i-1))
						if waitTime > s.maxRetryWait {
							waitTime = s.maxRetryWait
						}

						time.Sleep(time.Duration(waitTime) * time.Millisecond)

						err = s.deliver(hk, h, msg)
						if nil != err {
							continue
						}

						break
					}

					if nil != err {
						errs = append(errs, err)
					}

					wg.Done()
				}(hk, h)
			}

			wg.Wait()

			var err error
			var lvl zerolog.Level

			lvl = zerolog.DebugLevel

			if len(errs) > 0 {
				err = ErrHandlerOneOrMoreFail
				lvl = zerolog.ErrorLevel

				if len(errs) == len(handlers) {
					err = ErrHandlerAllFail
				}
			}

			if nil != s.metrics {
				go s.metrics.record(msg.event, err)
			}

			log.WithLevel(lvl).
				Err(err).
				Errs("errors", errs).
				Str("id", msg.event.ID()).
				Str("type", msg.event.Type()).
				Str("specversion", msg.event.SpecVersion()).
				Str("source", msg.event.Source()).
				Str("subject", msg.event.Subject()).
				Time("time", msg.event.Time()).
				Msg("Event process")

		case <-s.context.Done():
			break InfiniteLoop
		}
	}
}

func (s *Singularity) deliver(handlerkey string, h handler.Handler, msg message) error {
	var (
		err error
		lvl zerolog.Level
	)

	defer func() {
		rec := recover()

		log.WithLevel(lvl).
			Err(err).
			Interface("recover", rec).
			Str("id", msg.event.ID()).
			Str("type", msg.event.Type()).
			Str("specversion", msg.event.SpecVersion()).
			Str("source", msg.event.Source()).
			Str("subject", msg.event.Subject()).
			Str("handler", handlerkey).
			Str("output", h.Output().(asyncOutput).Ref()).
			Time("time", msg.event.Time()).
			Msg("Event handler deliver")
	}()

	lvl = zerolog.DebugLevel

	if s.validation {
		validator := s.findHandlerValidator(handlerkey)
		if nil != validator {
			match := validator.Match(msg.event)

			if false == match {
				return nil
			}
		}
	}

	err = h.Handle(msg.context, msg.event)
	if nil != err {
		lvl = zerolog.ErrorLevel
	}

	return err
}

func (s *Singularity) findHandlerValidator(handlerkey string) validator.Validator {
	if nil == s.validators[handlerkey] {
		return nil
	}

	return s.validators[handlerkey]
}

func (s *Singularity) Receiver() func(context.Context, cloudevents.Event) error {
	return s.queue
}

func (s *Singularity) SetOption(opt Option) {
	opt(s)
}

func (s *Singularity) Metrics() *metricServer {
	return s.metrics
}

func (s *Singularity) Close() {
	s.cancel()
}

func New(name string, opts ...Option) (*Singularity, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Singularity{
		name:       name,
		mutex:      &sync.RWMutex{},
		handlers:   map[string]handler.Handler{},
		validators: map[string]validator.Validator{},
		context:    ctx,
		cancel:     cancel,
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)

			if nil != s.context.Err() {
				break
			}

			if nil == s.pending {
				continue
			}

			log.WithLevel(level(len(s.pending), s.backlog)).
				Str("backlog", fmt.Sprintf("%d/%d", len(s.pending), s.backlog)).
				Msg("Pending events")
		}
	}()

	for _, opt := range opts {
		opt(s)
	}

	go s.process()

	return s, nil
}

func level(pending, size int) zerolog.Level {
	percent := pending * 100 / size

	// case >= 90%
	if percent >= 90 {
		return zerolog.WarnLevel
	}

	// case >= 70%
	if percent >= 70 {
		return zerolog.InfoLevel
	}

	// case < 70%
	return zerolog.DebugLevel
}

func rateinc(x, y float64) int {
	return int(math.Pow(x, y))
}
