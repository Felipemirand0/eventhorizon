package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EventHorizon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EventHorizonSpec `json:"spec"`
}

type EventHorizonSpec struct {
	Labels    map[string]string `json:"labels,omitempty"`
	Queue     *Queue            `json:"queue,omitempty"`
	Metrics   *Metrics          `json:"metrics,omitempty"`
	Transport *Transport        `json:"transport,omitempty"`
	Encoder   *Encoder          `json:"encoder"`
	Output    *Output           `json:"output,omitempty"`
	Validator *Validator        `json:"validator,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EventHorizonList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EventHorizon `json:"items"`
}

type Queue struct {
	Backlog      int `json:"backlog"`
	MaxRetry     int `json:"maxRetry"`
	RetryWait    int `json:"retryWait"`
	MaxRetryWait int `json:"maxRetryWait"`
}

type Metrics struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

type Transport struct {
	Type string         `json:"type"`
	HTTP *HTTPTransport `json:"http"`
	NATS *NATSTransport `json:"nats"`
}

type HTTPTransport struct {
	Port            int  `json:"port"`
	UseStatusCodeOK bool `json:"useStatusCodeOK"`
}

type NATSTransport struct {
	Server  string `json:"server"`
	Subject string `json:"subject"`
}

type Encoder struct {
	Type string `json:"type"`
}

type Output struct {
	Type    string         `json:"type"`
	Fluentd *OutputFluentd `json:"fluentd"`
}

type OutputFluentd struct {
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

type Validator struct {
	AllowedTypes   []string `json:"allowedTypes"`
	AllowedSources []string `json:"allowedSources"`
}
