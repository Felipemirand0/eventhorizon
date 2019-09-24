package validator

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

type Validator interface {
	Name() string
	Match(cloudevents.Event) bool
}
