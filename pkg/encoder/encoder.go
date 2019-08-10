package encoder

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go"
)

type Encoder interface {
	Kind() string
	Encode(context.Context, cloudevents.Event) (interface{}, error)
}
