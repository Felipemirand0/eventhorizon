package controller

import (
	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/singularity"

	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventsclient "github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) SyncSingularity(e *v1alpha1.Singularity) error {
	key, err := cache.MetaNamespaceKeyFunc(e)
	if nil != err {
		return err
	}

	if c.name != key {
		return ErrNameMismatch
	}

	if nil != c.singularity {
		return ErrAlreadyRunning
	}

	opts := []singularity.Option{}

	if e.Spec.Validation {
		opts = append(opts, singularity.EnableValidation())
	}

	c.singularity, err = singularity.New(key, opts...)
	if nil != err {
		return err
	}

	var transport transport.Transport

	switch e.Spec.Transport.Name {
	case "http":
		transport, err = cloudeventshttp.New(
			cloudeventshttp.WithPort(e.Spec.Transport.HTTP.Port),
		)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Set transport")

			return err
		}

		log.Info().
			Str("name", "http").
			Int("port", e.Spec.Transport.HTTP.Port).
			Msg("Set transport")

	case "nats":
		transport, err = cloudeventsnats.New(
			e.Spec.Transport.NATS.Server,
			e.Spec.Transport.NATS.Subject,
		)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Set transport")
		}

		log.Info().
			Str("name", "nats").
			Str("server", e.Spec.Transport.NATS.Server).
			Str("subject", e.Spec.Transport.NATS.Subject).
			Msg("Set transport")

	default:
		log.Error().
			Err(ErrUnknownTransport).
			Strs("available", []string{"nats", "http"}).
			Str("set", e.Spec.Transport.Name).
			Msg("Set transport")

		return ErrUnknownTransport
	}

	cli, err := cloudevents.NewClient(transport)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Create client")

		return err
	}

	go func(cli cloudeventsclient.Client) {
		var err error
		var lvl zerolog.Level

		defer func() {
			rec := recover()

			log.WithLevel(lvl).
				Err(err).
				Interface("recover", rec).
				Msg("Stop client")
		}()

		lvl = zerolog.InfoLevel

		err = cli.StartReceiver(c.context, c.singularity.Receiver())
		if nil != err {
			lvl = zerolog.ErrorLevel
		}
	}(cli)

	log.Info().
		Str("name", key).
		Str("transport", e.Spec.Transport.Name).
		Bool("validation", e.Spec.Validation).
		Msg("Set singularity")

	return nil
}
