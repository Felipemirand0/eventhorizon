package controller

import (
	"context"
	"time"

	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/output"
)

type asyncOutput struct {
	output     output.Output
	reference  string
	controller *Controller
	context    context.Context
	cancel     context.CancelFunc
}

func (o *asyncOutput) Send(ctx context.Context, e interface{}) error {
	if nil == o.output {
		return ErrNoOutputRegistered
	}

	return o.output.Send(ctx, e)
}

func (o *asyncOutput) Close() error {
	o.cancel()

	if nil == o.output {
		return nil
	}

	return o.output.Close()
}

func (o *asyncOutput) Ref() string {
	return o.reference
}

func (o *asyncOutput) watch() {
	for {
		time.Sleep(3 * time.Second)

		if nil != o.context.Err() {
			break
		}

		if nil != o.output {
			break
		}

		if nil == o.controller.outputs[o.reference] {
			continue
		}

		o.output = o.controller.outputs[o.reference]

		break
	}
}

func newAsyncOutput(ref string, c *Controller) *asyncOutput {
	ctx, cancel := context.WithCancel(context.Background())

	o := asyncOutput{
		reference:  ref,
		controller: c,
		context:    ctx,
		cancel:     cancel,
	}

	go o.watch()

	return &o
}
