package singularity

import (
	"context"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
)

type Receiver struct {
	cli             client.Client
	useStatusCodeOK bool
}

func (r *Receiver) SetStatusCodeOK(set bool) {
	r.useStatusCodeOK = set
}

func (r *Receiver) Receive(ctx context.Context, e cloudevents.Event, er *cloudevents.EventResponse) error {
	var err error

	_, ok := r.cli.(transport.Receiver)

	if ok {
		err = r.cli.(transport.Receiver).Receive(ctx, e, er)
	}

	if r.useStatusCodeOK {
		er.Status = http.StatusOK
	}

	return err
}

func NewReceiver(cli client.Client) *Receiver {
	return &Receiver{
		cli: cli,
	}
}
