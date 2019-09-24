package handler

import (
	"context"

	"acesso.io/eventhorizon/pkg/encoder"
	"acesso.io/eventhorizon/pkg/output"

	cloudevents "github.com/cloudevents/sdk-go"
)

type Basic struct {
	output  output.Output
	encoder encoder.Encoder
	labels  map[string]string
}

func (h *Basic) Handle(ctx context.Context, event cloudevents.Event) error {
	var (
		data interface{}
		err  error
	)

	if len(h.labels) > 0 {
		event.SetExtension("labels", h.labels)
	}

	if nil != h.encoder {
		data, err = h.encoder.Encode(ctx, event)
		if nil != err {
			return err
		}
	}

	return h.output.Send(ctx, data)
}

func (r *Basic) Close() error {
	return nil
}

func (r *Basic) Output() output.Output {
	return r.output
}

func NewBasic(out output.Output, enc encoder.Encoder, lab map[string]string) (*Basic, error) {
	h := Basic{
		output:  out,
		encoder: enc,
		labels:  lab,
	}

	return &h, nil
}
