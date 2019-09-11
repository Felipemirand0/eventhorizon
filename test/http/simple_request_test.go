package http_test

import (
	"context"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
)

func TestHTTPClient(t *testing.T) {
	trans, err := cloudeventshttp.New(
		cloudeventshttp.WithEncoding(cloudeventshttp.BinaryV03),
		cloudeventshttp.WithTarget("http://localhost:1257/"),
	)
	if err != nil {
		panic("failed to create transport: " + err.Error())
	}

	c, err := cloudevents.NewClient(trans)
	if err != nil {
		panic("failed to create client: " + err.Error())
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)

	event.SetSubject("MyMethod.MyAction")
	event.SetType("io.request.rpc")
	event.SetSource("myapp")
	event.SetID("TestHTTPClient")
	event.SetTime(time.Now())
	event.SetExtension("custom-a", "Foo")
	event.SetExtension("custom-b", "Bar")
	event.SetData(map[string]interface{}{
		"testing": "test",
	})

	_, _, err = c.Send(context.Background(), event)
	if err != nil {
		t.Errorf("failed delivery event: %v", err)
	}
}
