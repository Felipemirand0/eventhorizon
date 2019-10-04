module acesso.io/eventhorizon

go 1.12

require (
	cloud.google.com/go v0.45.1 // indirect
	github.com/Azure/go-autorest/autorest v0.9.1 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/cloudevents/sdk-go v0.9.2
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/fluent/fluent-logger-golang v1.4.0
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/nats-io/nats-server/v2 v2.0.4 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	github.com/rs/zerolog v1.15.0
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/vrischmann/envconfig v1.2.0
	go.opencensus.io v0.22.1 // indirect
	golang.org/x/crypto v0.0.0-20190909091759-094676da4a83 // indirect
	golang.org/x/exp v0.0.0-20190829153037-c13cbed26979 // indirect
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b // indirect
	golang.org/x/sys v0.0.0-20190909082730-f460065e899a // indirect
	golang.org/x/tools v0.0.0-20190909214602-067311248421 // indirect
	gonum.org/v1/gonum v0.0.0-20190908220844-1d8f8b2ee4ce // indirect
	google.golang.org/appengine v1.6.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190905160310-fb749d2f1064
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.0.0-20190831074504-732c9ca86353
	k8s.io/gengo v0.0.0-20190907103519-ebc107f98eab // indirect
	k8s.io/klog v0.4.0
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // indirect
	k8s.io/utils v0.0.0-20190907131718-3d4f5b7dea0b // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.3+incompatible
	k8s.io/api => k8s.io/api v0.0.0-20190620084959-7cf5895f2711
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190620085212-47dc9a115b18
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190620085130-185d68e6e6ea
)
