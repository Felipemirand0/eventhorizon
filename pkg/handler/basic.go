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
	if nil == h.encoder {
		return nil
	}

	if len(h.labels) > 0 {
		event.SetExtension("labels", h.labels)
	}

	data, err := h.encoder.Encode(ctx, event)
	if nil != err {
		return err
	}

	return h.output.Send(ctx, data)
}

func (r *Basic) Close() error {
	return nil
}

func NewBasic(out output.Output, enc encoder.Encoder, labels map[string]string) (*Basic, error) {
	h := Basic{
		output:  out,
		encoder: enc,
		labels:  labels,
	}

	return &h, nil
}
