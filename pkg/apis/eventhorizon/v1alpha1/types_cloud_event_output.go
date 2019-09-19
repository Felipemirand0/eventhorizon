package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventOutput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudEventOutputSpec `json:"spec"`
}

type CloudEventOutputSpec struct {
	Type    string                  `json:"type"`
	Fluentd CloudEventOutputFluentd `json:"fluentd"`
}

type CloudEventOutputFluentd struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	SocketPath         string `json:"socketPath"`
	Network            string `json:"network"`
	Timeout            string `json:"timeout"`
	WriteTimeout       string `json:"writeTimeout"`
	BufferLimit        int    `json:"bufferLimit"`
	RetryWait          int    `json:"retryWait"`
	MaxRetryWait       int    `json:"maxRetryWait"`
	MaxRetry           int    `json:"maxRetry"`
	Async              bool   `json:"async"`
	SubSecondPrecision bool   `json:"subSecondPrecision"`
	RequestAck         bool   `json:"requestAck"`
	TagPrefix          string `json:"tagPrefix"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CloudEventOutputList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CloudEventOutput `json:"items"`
}
