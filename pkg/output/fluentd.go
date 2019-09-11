package output

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Fluentd struct {
	cli *fluent.Fluent
}

func (o *Fluentd) Send(ctx context.Context, event interface{}) error {
	err := o.cli.Post("eventhorizon", event)
	if nil != err {
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Str("output", "fluentd").
			Msg("failed to delivery message")
	}

	return err
}

func (o *Fluentd) Close() error {
	return o.cli.Close()
}

func NewFluentd(cfg fluent.Config) (*Fluentd, error) {
	var err error

	o := Fluentd{}

	waitForFluentd(cfg)

	o.cli, err = fluent.New(cfg)
	if err != nil {
		log.WithLevel(zerolog.WarnLevel).
			Err(err).
			Str("output", "fluentd").
			Msg("failed to connect")

		return nil, err
	}

	return &o, nil
}

func waitForFluentd(cfg fluent.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for {
		if ctx.Err() != nil {
			return errors.New("unable to connect to fluentd")
		}

		switch cfg.FluentNetwork {
		case "unix":
			if _, err := os.Stat(cfg.FluentSocketPath); err != nil {
				if os.IsNotExist(err) {
					continue
				}

				return err
			}

		case "tcp":
			_, err := net.Dial("tcp", net.JoinHostPort(cfg.FluentHost, strconv.Itoa(cfg.FluentPort)))

			if nil != err {
				continue
			}
		}

		break
	}

	return nil
}
