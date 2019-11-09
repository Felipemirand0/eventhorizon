package errors

import "errors"

var (
	ErrNameMismatch     = errors.New("Resource did not match instance name")
	ErrUnknownTransport = errors.New("Unknow Tranport")
	ErrUnknownOutput    = errors.New("Unknow Output")
	ErrUnknownEncoder   = errors.New("Unknow Encoder")
	ErrWaitCacheSync    = errors.New("Failed waiting for informer caches to sync")
)
