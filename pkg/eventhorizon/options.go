package eventhorizon

import (
	"time"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"
	"acesso.io/eventhorizon/pkg/encoder"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/metrics"
	"acesso.io/eventhorizon/pkg/output"
	"acesso.io/eventhorizon/pkg/validator"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Option func(*EventHorizon) error

func SetQueue(tmp v1alpha2.Queue) Option {
	return func(s *EventHorizon) error {
		var (
			Backlog      = 8192
			MaxRetry     = 25
			RetryWait    = 500
			MaxRetryWait = 65000
		)

		if tmp.Backlog > 0 {
			Backlog = tmp.Backlog
		}

		if tmp.MaxRetry > 0 {
			MaxRetry = tmp.MaxRetry
		}

		if tmp.RetryWait > 0 {
			RetryWait = tmp.RetryWait
		}

		if tmp.MaxRetryWait > 0 {
			MaxRetryWait = tmp.MaxRetryWait
		}

		s.backlog = Backlog
		s.maxRetry = MaxRetry
		s.retryWait = RetryWait
		s.maxRetryWait = MaxRetryWait

		return nil
	}
}

func SetTransport(tmp v1alpha2.Transport) Option {
	return func(s *EventHorizon) error {
		var t transport.Transport
		var err error
		var fields map[string]interface{}
		var level zerolog.Level = zerolog.InfoLevel

		defer func() {
			log.WithLevel(level).
				Err(err).
				Str("type", tmp.Type).
				Fields(fields).
				Msg("Set Transport")
		}()

		switch tmp.Type {
		case "http":
			var (
				Port        = 1257
				UseStatusOK = false
			)

			if nil != tmp.HTTP {
				if tmp.HTTP.Port > 0 {
					Port = tmp.HTTP.Port
				}

				if tmp.HTTP.UseStatusCodeOK {
					UseStatusOK = true
				}
			}

			t, err = cloudeventshttp.New(cloudeventshttp.WithPort(Port), WithHTTPServerMux(s.healthcheck))
			if err != nil {
				break
			}

			if UseStatusOK {
				s.useHTTPStatusOK = true
			}

			fields = map[string]interface{}{
				"port":            tmp.HTTP.Port,
				"useStatusCodeOK": tmp.HTTP.UseStatusCodeOK,
			}

		case "nats":
			var (
				Server  = "localhost:4222"
				Subject = "eventhorizon"
			)

			if nil != tmp.NATS {
				if tmp.NATS.Server != "" {
					Server = tmp.NATS.Server
				}

				if tmp.NATS.Subject != "" {
					Subject = tmp.NATS.Subject
				}
			}

			t, err = cloudeventsnats.New(Server, Subject)
			if err != nil {
				break
			}

			fields = map[string]interface{}{
				"server":  tmp.NATS.Server,
				"subject": tmp.NATS.Subject,
			}

		default:
			err = ErrUnknownTransport

			fields = map[string]interface{}{
				"available": []string{"nats", "http"},
			}
		}

		if nil != err {
			level = zerolog.ErrorLevel
		}

		s.transport = t

		return err
	}
}

func SetEncoder(tmp v1alpha2.Encoder) Option {
	return func(s *EventHorizon) error {
		var e encoder.Encoder
		var err error
		var fields map[string]interface{}
		var level zerolog.Level = zerolog.InfoLevel

		defer func() {
			log.WithLevel(level).
				Err(err).
				Str("type", tmp.Type).
				Fields(fields).
				Msg("Set Encoder")
		}()

		switch tmp.Type {
		case "map":
			e, err = encoder.NewMap()

		default:
			err = ErrUnknownEncoder

			fields = map[string]interface{}{
				"available": []string{"map"},
			}
		}

		if nil != err {
			level = zerolog.ErrorLevel
		}

		s.encoder = e

		return err
	}
}

func SetOutput(tmp v1alpha2.Output) Option {
	return func(s *EventHorizon) error {
		var o output.Output
		var err error
		var fields map[string]interface{}
		var level zerolog.Level = zerolog.InfoLevel

		defer func() {
			log.WithLevel(level).
				Err(err).
				Str("type", tmp.Type).
				Fields(fields).
				Msg("Set Output")
		}()

		switch tmp.Type {
		case "fluentd":
			var (
				timeout      time.Duration
				writeTimeout time.Duration
			)

			timeout, err = time.ParseDuration(tmp.Fluentd.Timeout)
			if err != nil {
				break
			}

			writeTimeout, err = time.ParseDuration(tmp.Fluentd.WriteTimeout)
			if err != nil {
				break
			}

			fluentcfg := fluent.Config{
				Timeout:            timeout,
				WriteTimeout:       writeTimeout,
				MaxRetry:           tmp.Fluentd.MaxRetry,
				FluentPort:         tmp.Fluentd.Port,
				FluentHost:         tmp.Fluentd.Host,
				FluentNetwork:      tmp.Fluentd.Network,
				FluentSocketPath:   tmp.Fluentd.SocketPath,
				BufferLimit:        tmp.Fluentd.BufferLimit,
				RetryWait:          tmp.Fluentd.RetryWait,
				MaxRetryWait:       tmp.Fluentd.MaxRetryWait,
				TagPrefix:          tmp.Fluentd.TagPrefix,
				Async:              tmp.Fluentd.Async,
				SubSecondPrecision: tmp.Fluentd.SubSecondPrecision,
				RequestAck:         tmp.Fluentd.RequestAck,
			}

			o, err = output.NewFluentd(fluentcfg)

			fields = map[string]interface{}{
				"network": tmp.Fluentd.Network,
				"host":    tmp.Fluentd.Host,
				"port":    tmp.Fluentd.Port,
				"socket":  tmp.Fluentd.SocketPath,
			}

		case "stdout":
			o, err = output.NewStdout()

		default:
			err = ErrUnknownOutput

			fields = map[string]interface{}{
				"available": []string{"fluentd", "stdout"},
			}
		}

		if nil != err {
			level = zerolog.ErrorLevel
		}

		s.output = o

		return err
	}
}

func WithValidator(validator validator.Validator) Option {
	return func(s *EventHorizon) error {
		s.validator = validator

		return nil
	}
}

func WithMetrics(metrics metrics.Metrics) Option {
	return func(s *EventHorizon) error {
		s.metrics = metrics

		return nil
	}
}

func WithLabels(labels map[string]string) Option {
	return func(s *EventHorizon) error {
		s.labels = labels

		return nil
	}
}
