package errors

import "errors"

var (
	ErrNoTransport          = errors.New("no transport set")
	ErrUnknownTransport     = errors.New("unknow tranport")
	ErrUnknownOutput        = errors.New("unknow output")
	ErrAlreadyRunning       = errors.New("service already running, live reload not implemented yet, restart manually")
	ErrNameMismatch         = errors.New("resource did not match service name")
	ErrNoMatchingSubject    = errors.New("resource has no matching subjects")
	ErrHandlerOneOrMoreFail = errors.New("one or more handler failed")
	ErrHandlerAllFail       = errors.New("all handlers failed")
	ErrNoHandlerRegistered  = errors.New("no handler registered")
	ErrNoOutputRegistered   = errors.New("no output registered")
	ErrWaitCacheSync        = errors.New("failed waiting cache sync")
)
