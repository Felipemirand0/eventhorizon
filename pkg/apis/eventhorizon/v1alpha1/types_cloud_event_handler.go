package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventHandler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudEventHandlerSpec `json:"spec"`
}

type CloudEventHandlerSpec struct {
	Subjects []string          `json:"subjects"`
	Labels   map[string]string `json:"labels,omitempty"`
	Encoder  string            `json:"encoder"`
	Output   string            `json:"output"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventHandlerList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CloudEventHandler `json:"items"`
}
