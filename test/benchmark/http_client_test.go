package eventhorizon

import (
	"context"
	"fmt"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
)

func BenchmarkHTTPClient(b *testing.B) {
	trans, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("http://localhost:1257/"),
		cloudevents.WithEncoding(cloudevents.HTTPBinaryV03),
	)
	if err != nil {
		panic("failed to create transport: " + err.Error())
	}

	c, err := cloudevents.NewClient(trans)
	if err != nil {
		panic("failed to create client: " + err.Error())
	}

	for n := 0; n < b.N; n++ {
		event := cloudevents.NewEvent()

		event.SetSubject("MyMethod.MyAction")
		event.SetType("io.request.rpc")
		event.SetSource("myapp")
		event.SetID(fmt.Sprintf("BenchmarkHTTPClient#%d", n))
		event.SetTime(time.Now())
		event.SetExtension("custom-a", "Foo")
		event.SetExtension("custom-b", "Bar")
		event.SetData(map[string]interface{}{
			"testing": "test",
		})

		_, _, err = c.Send(context.Background(), event)
		if err != nil {
			b.Errorf("failed delivery event: %v", err)
		}
	}
}
