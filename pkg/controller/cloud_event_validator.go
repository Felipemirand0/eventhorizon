package controller

import (
	"errors"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/validator"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) SyncCloudEventValidator(e *v1alpha1.CloudEventValidator) error {
	key, err := cache.MetaNamespaceKeyFunc(e)
	if nil != err {
		return err
	}

	if nil != c.validators[key] {
		return ErrAlreadyRunning
	}

	var val validator.Validator

	if nil == val {
		log.Warn().
			Err(errors.New("feature not implemented yet")).
			Str("name", key).
			Strs("handlers", e.Spec.Handlers).
			Msg("Set validator")

		return nil
	}

	c.mutex.Lock()
	c.validators[key] = val
	c.mutex.Unlock()

	return nil
}
