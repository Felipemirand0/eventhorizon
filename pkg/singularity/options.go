package singularity

import (
	"acesso.io/eventhorizon/pkg/handler"
	"acesso.io/eventhorizon/pkg/validator"
)

type Option func(*Singularity)

func SetHandler(key string, handler handler.Handler) Option {
	return func(s *Singularity) {
		s.mutex.Lock()
		s.handlers[key] = handler
		s.mutex.Unlock()
	}
}

type RetryOptions struct {
	Backlog      int
	MaxRetry     int
	RetryWait    int
	MaxRetryWait int
}

func SetRetryOptions(opts RetryOptions) Option {
	return func(s *Singularity) {
		s.pending = make(chan message, opts.Backlog)
		s.backlog = opts.Backlog
		s.maxRetry = opts.MaxRetry
		s.retryWait = opts.RetryWait
		s.maxRetryWait = opts.MaxRetryWait
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
