package validator

import (
	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	cloudevents "github.com/cloudevents/sdk-go"
)

func NewBasic(e *v1alpha1.CloudEventValidator) Validator {
	return &Basic{
		entity: e,
	}
}

type Basic struct {
	entity *v1alpha1.CloudEventValidator
}

func (v Basic) Name() string {
	return v.entity.Name
}

func (v Basic) Match(event cloudevents.Event) bool {
	// newEvent := cloudevents.NewEvent()
	// *event = newEvent

	return true
}
