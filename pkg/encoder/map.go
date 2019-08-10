package encoder

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go"
)

func NewMap() Encoder {
	return &Map{}
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
		"extensions":  event.Extensions(),
		"type":        event.Type(),
		"time":        event.Time().Unix(),
		"data":        data,
	}

	return val, nil
}
