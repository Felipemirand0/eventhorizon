package controller

import (
	"time"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/output"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) SyncCloudEventOutput(e *v1alpha1.CloudEventOutput) error {
	key, err := cache.MetaNamespaceKeyFunc(e)
	if nil != err {
		return err
	}

	var out output.Output

	switch e.Spec.Type {
	case "fluentd":
		timeout, err := time.ParseDuration(e.Spec.Fluentd.Timeout)
		if err != nil {
			return err
		}

		writeTimeout, err := time.ParseDuration(e.Spec.Fluentd.WriteTimeout)
		if err != nil {
			return err
		}

		fluentcfg := fluent.Config{
			Timeout:            timeout,
			WriteTimeout:       writeTimeout,
			MaxRetry:           e.Spec.Fluentd.MaxRetry,
			FluentPort:         e.Spec.Fluentd.Port,
			FluentHost:         e.Spec.Fluentd.Host,
			FluentNetwork:      e.Spec.Fluentd.Network,
			FluentSocketPath:   e.Spec.Fluentd.SocketPath,
			BufferLimit:        e.Spec.Fluentd.BufferLimit,
			RetryWait:          e.Spec.Fluentd.RetryWait,
			MaxRetryWait:       e.Spec.Fluentd.MaxRetryWait,
			TagPrefix:          e.Spec.Fluentd.TagPrefix,
			Async:              e.Spec.Fluentd.Async,
			SubSecondPrecision: e.Spec.Fluentd.SubSecondPrecision,
		}

		out, err = output.NewFluentd(fluentcfg)
		if nil != err {
			log.Error().
				Err(err).
				Msg("Failing to create output")

			return err
		}

	case "stdout":
		out, err = output.NewStdout()
		if nil != err {
			log.Error().
				Err(err).
				Msg("Failing to create output")

			return err
		}

	default:
		log.Error().
			Err(ErrUnknownOutput).
			Strs("available", []string{"fluentd", "stdout"}).
			Str("set", key).
			Msg("Unknow output")

		return ErrUnknownOutput
	}

	c.mutex.Lock()
	c.outputs[key] = out
	c.mutex.Unlock()

	go func() {
		<-c.context.Done()
		out.Close()
	}()

	log.Info().
		Str("name", key).
		Str("type", e.Spec.Type).
		Msg("Set output")

	return nil
}
