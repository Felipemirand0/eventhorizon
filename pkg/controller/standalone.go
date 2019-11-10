package controller

import (
	"context"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

type StandaloneController struct {
	Controller
}

func (c *StandaloneController) Run() error {
	log.Info().
		Msg("Starting EventHorizon")

	err := c.eventhorizon.Start()

	c.eventhorizon.Close()

	var level zerolog.Level = zerolog.InfoLevel

	if nil != err {
		level = zerolog.ErrorLevel
	}

	log.WithLevel(level).
		Err(err).
		Msg("Stopping EventHorizon")

	return nil
}

func NewStandalone(ctx context.Context, name string, r *v1alpha2.EventHorizon) *StandaloneController {
	c := &StandaloneController{
		Controller: Controller{
			name:    name,
			context: ctx,
		},
	}

	e, err := c.syncEventHorizon(r)
	if nil != err {
		key, _ := cache.MetaNamespaceKeyFunc(r)

		log.Fatal().
			Err(err).
			Str("name", name).
			Str("key", key).
			Msg("Failing resource to load")
	}

	c.eventhorizon = e

	return c
}
