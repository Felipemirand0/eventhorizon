package singularity

import (
	"acesso.io/eventhorizon/pkg/handler"
	"acesso.io/eventhorizon/pkg/validator"
)

type Option func(*Singularity)

func SetHandler(key string, handler handler.Handler) Option {
	return func(e *Singularity) {
		e.mutex.Lock()
		e.handlers[key] = handler
		e.mutex.Unlock()
	}
}

func SetHandlerValidator(handlerKey string, validator validator.Validator) Option {
	return func(s *Singularity) {
		s.mutex.Lock()
		s.validators[handlerKey] = validator
		s.mutex.Unlock()
	}
}

func EnableValidation() Option {
	return func(s *Singularity) {
		s.validation = true
	}
}

func EnableMetrics(path string, port int) Option {
	return func(s *Singularity) {
		s.metrics = newMetricsServer(path, port)
	}
}
