package eventhorizon

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"acesso.io/eventhorizon/pkg/encoder"
	"acesso.io/eventhorizon/pkg/handler"
	"acesso.io/eventhorizon/pkg/metrics"
	"acesso.io/eventhorizon/pkg/output"
	"acesso.io/eventhorizon/pkg/validator"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventHorizon struct {
	context     context.Context
	cancel      context.CancelFunc
	mutex       *sync.RWMutex
	healthcheck *healthcheck

	backlog      int
	maxRetry     int
	retryWait    int
	maxRetryWait int

	metrics   metrics.Metrics
	transport transport.Transport
	encoder   encoder.Encoder
	output    output.Output
	validator validator.Validator
	labels    map[string]string

	handler handler.Handler

	pending chan message

	useHTTPStatusOK bool
}

type message struct {
	context context.Context
	event   cloudevents.Event
}

func (s *EventHorizon) Start() error {
	s.pending = make(chan message, s.backlog)

	go func() {
		for {
			time.Sleep(10 * time.Second)

			if nil != s.context.Err() {
				break
			}

			log.WithLevel(level(len(s.pending), s.backlog)).
				Str("backlog", fmt.Sprintf("%d/%d", len(s.pending), s.backlog)).
				Msg("Pending events")
		}
	}()

	go s.process()

	if nil != s.metrics {
		go s.metrics.Listen(s.context)
	}

	cli, err := cloudevents.NewClient(s.transport)
	if err != nil {
		return err
	}

	if s.useHTTPStatusOK {
		r := newReceiver(cli)
		r.UseStatusCodeOK(true)

		// overwrite receiver
		s.transport.SetReceiver(r)
	}

	log.Info().
		Msg("Start EventHorizon")

	return cli.StartReceiver(s.context, s.Receiver())
}

func (s *EventHorizon) queue(ctx context.Context, e cloudevents.Event) error {
	s.pending <- message{
		context: ctx,
		event:   e,
	}

	return nil
}

func (s *EventHorizon) process() {
InfiniteLoop:
	for {
		select {
		case msg := <-s.pending:
			err := s.deliver(msg)

			if nil != s.metrics {
				go s.metrics.Record(msg.event, err)
			}

			if nil == err {
				s.healthcheck.failing = true
				s.healthcheck.reason = err
			} else {
				s.healthcheck.failing = false
				s.healthcheck.reason = nil
			}

		case <-s.context.Done():
			break InfiniteLoop
		}
	}
}

func (s *EventHorizon) deliver(msg message) error {
	var (
		errs  []error
		err   error
		match bool          = true
		level zerolog.Level = zerolog.DebugLevel
	)

	for i := 0; i < s.maxRetry; i++ {
		err = s.handler.Handle(msg.context, msg.event)

		if nil != err {
			errs = append(errs, err)

			waitTime := s.retryWait * rateinc(1.5, float64(i-1))
			if waitTime > s.maxRetryWait {
				waitTime = s.maxRetryWait
			}

			time.Sleep(time.Duration(waitTime) * time.Millisecond)
			continue
		}

		break
	}

	if nil == err && len(errs) > 0 {
		level = zerolog.WarnLevel
	}

	if nil != err {
		level = zerolog.ErrorLevel
	}

	log.WithLevel(level).
		Err(err).
		Errs("errors", errs).
		Bool("match", match).
		Str("id", msg.event.ID()).
		Str("type", msg.event.Type()).
		Str("specversion", msg.event.SpecVersion()).
		Str("source", msg.event.Source()).
		Str("subject", msg.event.Subject()).
		Time("time", msg.event.Time()).
		Msg("Event handler deliver")

	return err
}

func (s *EventHorizon) Receiver() func(context.Context, cloudevents.Event) error {
	return s.queue
}

func (s *EventHorizon) SetOption(opt Option) {
	opt(s)
}

func (s *EventHorizon) Close() {
	log.Info().
		Msg("Stop EventHorizon")

	s.cancel()
}

func New(opts ...Option) (*EventHorizon, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &EventHorizon{
		context:     ctx,
		cancel:      cancel,
		mutex:       &sync.RWMutex{},
		healthcheck: &healthcheck{},
	}

	for _, opt := range opts {
		err := opt(s)

		if nil != err {
			return nil, err
		}
	}

	var err error

	s.handler, err = handler.NewBasic(s.output, s.encoder, s.labels)

	if nil != err {
		return nil, err
	}

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
