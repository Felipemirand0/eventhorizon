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

func (c *Controller) syncEventHorizon(r *v1alpha2.EventHorizon) (*eventhorizon.EventHorizon, error) {
	key, err := cache.MetaNamespaceKeyFunc(r)
	if nil != err {
		return nil, err
	}

	if key != c.name {
		return nil, ErrNameMismatch
	}

	que := defaultQueue
	if nil != r.Spec.Queue {
		que = *r.Spec.Queue
	}

	tra := defaultTranport
	if nil != r.Spec.Transport {
		tra = *r.Spec.Transport
	}

	enc := defaultEncoder
	if nil != r.Spec.Encoder {
		enc = *r.Spec.Encoder
	}

	out := defaultOutput
	if nil != r.Spec.Output {
		out = *r.Spec.Output
	}

	opts := []eventhorizon.Option{
		eventhorizon.SetQueue(que),
		eventhorizon.SetTransport(tra),
		eventhorizon.SetEncoder(enc),
		eventhorizon.SetOutput(out),
	}

	if nil != r.Spec.Validator {
		val, err := validator.NewBasic(*r.Spec.Validator)
		if nil == err {
			opts = append(opts, eventhorizon.WithValidator(val))
		}
	}

	if nil != r.Spec.Metrics {
		met, err := metrics.NewBasic(*r.Spec.Metrics)
		if nil == err {
			opts = append(opts, eventhorizon.WithMetrics(met))
		}
	}

	if len(r.Spec.Labels) > 0 {
		opts = append(opts, eventhorizon.WithLabels(r.Spec.Labels))
	}

	e, err := eventhorizon.New(c.context, opts...)
	if nil != err {
		return nil, err
	}

	log.Info().
		Str("name", key).
		Str("transport", r.Spec.Transport.Type).
		Str("encoder", r.Spec.Encoder.Type).
		Str("output", r.Spec.Output.Type).
		Msg("Set EventHorizon")

	return e, nil
}
