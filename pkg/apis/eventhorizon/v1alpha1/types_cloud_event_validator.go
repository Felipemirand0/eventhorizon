package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudEventValidatorSpec `json:"spec"`
}

type CloudEventValidatorSpec struct {
	Handlers       []string `json:"handlers"`
	Priority       int      `json:"priority"`
	AllowedTypes   []string `json:"allowedTypes"`
	AllowedSources []string `json:"allowedSources"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CloudEventValidator `json:"items"`
}
