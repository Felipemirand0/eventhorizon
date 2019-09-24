package handler

import (
	"context"

	"acesso.io/eventhorizon/pkg/output"

	cloudevents "github.com/cloudevents/sdk-go"
)

type Handler interface {
	Handle(context.Context, cloudevents.Event) error
	Close() error
	Output() output.Output
}
