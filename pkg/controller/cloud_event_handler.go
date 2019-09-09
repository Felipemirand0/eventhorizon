package controller

import (
	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	"acesso.io/eventhorizon/pkg/encoder"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/handler"
	. "acesso.io/eventhorizon/pkg/helpers"
	"acesso.io/eventhorizon/pkg/output"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) SyncCloudEventHandler(e *v1alpha1.CloudEventHandler) error {
	if false == In(c.name, e.Spec.Subjects) {
		return ErrNoMatchingSubject
	}

	key, err := cache.MetaNamespaceKeyFunc(e)
	if nil != err {
		return err
	}

	var han handler.Handler
	var enc encoder.Encoder
	var out output.Output

	switch e.Spec.Encoder {
	case "map":
		enc = encoder.NewMap()
	}

	out = &outputWrapper{
		reference:  e.Spec.Output,
		controller: c,
	}

	han, err = handler.NewBasic(out, enc, e.Spec.Labels)
	if nil != err {
		log.Error().
			Err(err).
			Msg("Set handler")

		return err
	}

	c.mutex.Lock()
	c.handlers[key] = han
	c.mutex.Unlock()

	go func() {
		<-c.context.Done()
		han.Close()
	}()

	log.Info().
		Str("name", key).
		Str("output", e.Spec.Output).
		Str("encoder", e.Spec.Encoder).
		Strs("subjects", e.Spec.Subjects).
		Msg("Set handler")

	return nil
}
