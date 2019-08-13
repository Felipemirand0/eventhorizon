package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Singularity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SingularitySpec `json:"spec"`
}

type SingularitySpec struct {
	Transport  *Transport `json:"transport"`
	Validation bool       `json:"validation"`
	Metrics    *Metrics   `json:"metrics"`
}

type Transport struct {
	Name string         `json:"name"`
	HTTP *HTTPTransport `json:"http"`
	NATS *NATSTransport `json:"nats"`
}

type HTTPTransport struct {
	Port int `json:"port"`
}

type NATSTransport struct {
	Server  string `json:"server"`
	Subject string `json:"subject"`
}

type Metrics struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SingularityList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Singularity `json:"items"`
}
