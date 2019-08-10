package handler

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go"
)

type Handler interface {
	Handle(context.Context, cloudevents.Event) error
	Close() error
}
