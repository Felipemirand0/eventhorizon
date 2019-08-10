package validator

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

type Validator interface {
	Name() string
	Validate(cloudevents.Event) error
}
