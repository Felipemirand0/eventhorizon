package controller

import (
	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/eventhorizon"
	"acesso.io/eventhorizon/pkg/metrics"
	"acesso.io/eventhorizon/pkg/validator"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

var defaultQueue = v1alpha2.Queue{
	Backlog:      8192,
	MaxRetry:     25,
	RetryWait:    500,
	MaxRetryWait: 65000,
}

var defaultTranport = v1alpha2.Transport{
	Type: "http",
	HTTP: &v1alpha2.HTTPTransport{
		Port: 1257,
	},
}

var defaultOutput = v1alpha2.Output{
	Type: "stdout",
}

var defaultEncoder = v1alpha2.Encoder{
	Type: "map",
}

func (c *Controller) SyncEventHorizon(e *v1alpha2.EventHorizon) error {
	key, err := cache.MetaNamespaceKeyFunc(e)
	if nil != err {
		return err
	}

	if key != c.name {
		return ErrNameMismatch
	}

	que := defaultQueue
	if nil != e.Spec.Queue {
		que = *e.Spec.Queue
	}

	tra := defaultTranport
	if nil != e.Spec.Transport {
		tra = *e.Spec.Transport
	}

	enc := defaultEncoder
	if nil != e.Spec.Encoder {
		enc = *e.Spec.Encoder
	}

	out := defaultOutput
	if nil != e.Spec.Output {
		out = *e.Spec.Output
	}

	opts := []eventhorizon.Option{
		eventhorizon.SetQueue(que),
		eventhorizon.SetTransport(tra),
		eventhorizon.SetEncoder(enc),
		eventhorizon.SetOutput(out),
	}

	if nil != e.Spec.Validator {
		val, err := validator.NewBasic(*e.Spec.Validator)
		if nil == err {
			opts = append(opts, eventhorizon.WithValidator(val))
		}
	}

	if nil != e.Spec.Metrics {
		met, err := metrics.NewBasic(*e.Spec.Metrics)
		if nil == err {
			opts = append(opts, eventhorizon.WithMetrics(met))
		}
	}

	if len(e.Spec.Labels) > 0 {
		opts = append(opts, eventhorizon.WithLabels(e.Spec.Labels))
	}

	c.eventhorizon, err = eventhorizon.New(opts...)
	if nil != err {
		return err
	}

	log.Info().
		Str("name", key).
		Str("transport", e.Spec.Transport.Type).
		Str("encoder", e.Spec.Encoder.Type).
		Str("output", e.Spec.Output.Type).
		Msg("Set EventHorizon")

	return nil
}
