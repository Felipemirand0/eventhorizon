package validator

import (
	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"
	cloudevents "github.com/cloudevents/sdk-go"
)

func NewBasic(e v1alpha2.Validator) (Validator, error) {
	return &Basic{
		allowedTypes:   e.AllowedTypes,
		allowedSources: e.AllowedSources,
	}, nil
}

type Basic struct {
	allowedTypes   []string
	allowedSources []string
}

func (v Basic) Match(event cloudevents.Event) bool {
	return true
}
