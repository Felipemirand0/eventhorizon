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
}

func (h *Basic) Handle(ctx context.Context, event cloudevents.Event) error {
	var (
		data interface{}
		err  error
	)

	if nil != h.encoder {
		data, err = h.encoder.Encode(ctx, event)
		if nil != err {
			return err
		}
	}

	if nil == h.output {
		return nil
	}

	return h.output.Send(ctx, data)
}

func (r *Basic) Close() error {
	return nil
}

func NewBasic(out output.Output, enc encoder.Encoder) (*Basic, error) {
	h := Basic{
		output:  out,
		encoder: enc,
	}

	return &h, nil
}
