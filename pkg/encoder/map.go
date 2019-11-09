package encoder

import (
	"context"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
)

func NewMap() (Encoder, error) {
	return &Map{}, nil
}

type Map struct {
}

func (e Map) Kind() string {
	return "map"
}

func (e Map) Encode(ctx context.Context, event cloudevents.Event) (interface{}, error) {
	data := map[string]interface{}{}

	err := event.DataAs(&data)
	if nil != err {
		return nil, err
	}

	val := map[string]interface{}{
		"specversion": event.SpecVersion(),
		"id":          event.ID(),
		"subject":     event.Subject(),
		"source":      event.Source(),
		"type":        event.Type(),
		"time":        event.Time().Format(time.RFC3339),
		"extensions":  event.Extensions(),
		"data":        data,
	}

	return val, nil
}
