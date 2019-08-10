package output

import (
	"context"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Fluentd struct {
	cli *fluent.Fluent
}

func (o *Fluentd) Send(ctx context.Context, event interface{}) error {
	return o.cli.Post("eventhorizon", event)
}

func (o *Fluentd) Close() error {
	return o.cli.Close()
}

func NewFluentd(cfg fluent.Config) (*Fluentd, error) {
	var err error

	o := Fluentd{}

	o.cli, err = fluent.New(cfg)
	if err != nil {
		log.WithLevel(zerolog.WarnLevel).
			Err(err).
			Msg("Failing to connect to fluentd")

		return &o, err
	}

	return &o, nil
}
