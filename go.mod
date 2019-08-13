module acesso.io/eventhorizon

go 1.12

require (
	github.com/cloudevents/sdk-go v0.8.0
	github.com/fluent/fluent-logger-golang v1.4.0
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/rs/zerolog v1.14.3
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/vrischmann/envconfig v1.2.0
	k8s.io/api v0.0.0-20190808180749-077ce48e77da
	k8s.io/apimachinery v0.0.0-20190809020650-423f5d784010
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.0.0-00010101000000-000000000000
	k8s.io/klog v0.4.0
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
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
