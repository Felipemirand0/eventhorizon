package controller

import (
	"context"
	"fmt"
)

type outputWrapper struct {
	reference  string
	controller *Controller
}

func (o *outputWrapper) Send(ctx context.Context, e interface{}) error {
	if nil == o.controller.outputs[o.reference] {
		return fmt.Errorf("no registered output: %s", o.reference)
	}

	return o.controller.outputs[o.reference].Send(ctx, e)
}

func (o *outputWrapper) Close() error {
	if nil == o.controller.outputs[o.reference] {
		return nil
	}

	return o.controller.outputs[o.reference].Close()
}
