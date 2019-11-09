package validator

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

type Validator interface {
	Match(cloudevents.Event) bool
}
