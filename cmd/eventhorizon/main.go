package main

import (
	"io/ioutil"
	rawlog "log"
	"os"
	"time"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	"acesso.io/eventhorizon/pkg/controller"
	clientset "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	acessoschema "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	informers "acesso.io/eventhorizon/pkg/generated/informers/externalversions"
	"acesso.io/eventhorizon/pkg/signals"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

type conf struct {
	Mode        string `envconfig:"default=kubernetes"`
	Name        string `envconfig:"default=kube-system/eventhorizon"`
	Threadiness int    `envconfig:"default=1"`
	Standalone  struct {
		Config string `envconfig:"default=/opt/acesso/samples/standalone/stdout.yml"`
	}
	Kubernetes struct {
		InCluster bool `envconfig:"default=true"`
		// Path to a kubeconfig. Only required if out-of-cluster.
		KubeConfig string `envconfig:"optional"`
		// The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
		MasterURL string `envconfig:"optional"`
	}
	Logging struct {
		Enabled bool   `envconfig:"default=true"`
		Level   string `envconfig:"default=info"`
		Pretty  bool   `envconfig:"default=false"`
	}
}

var env conf

func init() {
	zerolog.TimestampFieldName = "ts"
}

func kubernetes() {
	stopCh := signals.SetupSignalHandler()

	var cfg *rest.Config
	var err error

	if env.Kubernetes.InCluster {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			klog.Fatalf("Error building config: %s", err.Error())
		}
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags(env.Kubernetes.MasterURL, env.Kubernetes.KubeConfig)
		if err != nil {
			klog.Fatalf("Error building config: %s", err.Error())
		}
	}

	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)

	c := controller.NewKubernetes(env.Name, env.Threadiness, client, informerFactory.Eventhorizon())

	informerFactory.Start(stopCh)

	if err = c.Run(stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func standalone() {
	stopCh := signals.SetupSignalHandler()

	acessoschema.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode

	var list *v1.List

	data, err := ioutil.ReadFile(env.Standalone.Config)
	if nil != err {
		log.Fatal().
			Err(err).
			Str("file", env.Standalone.Config).
			Msg("Reading configuration file")
	}

	obj, _, err := decode(data, nil, nil)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("file", env.Standalone.Config).
			Msg("Decoding configuration file")
	}

	switch obj.(type) {
	case *v1.List:
		list = obj.(*v1.List)

	default:
		log.Fatal().
			Msg("Configuration resource file must be of kind `List`")
	}

	var (
		singularity         *v1alpha1.Singularity
		cloudEventOutput    *v1alpha1.CloudEventOutput
		cloudEventHandler   *v1alpha1.CloudEventHandler
		cloudEventValidator *v1alpha1.CloudEventValidator
	)

	for _, item := range list.Items {
		obj, _, err := decode(item.Raw, nil, nil)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Decoding configuration file")
		}

		switch e := obj.(type) {
		case *v1alpha1.Singularity:
			singularity = e

		case *v1alpha1.CloudEventOutput:
			cloudEventOutput = e

		case *v1alpha1.CloudEventHandler:
			cloudEventHandler = e

		case *v1alpha1.CloudEventValidator:
			cloudEventValidator = e
		}
	}

	c := controller.NewStandalone(env.Name)

	err = c.SyncSingularity(singularity)
	if nil != err {
		resource, _ := cache.MetaNamespaceKeyFunc(singularity)

		log.Fatal().
			Err(err).
			Str("singularity", env.Name).
			Str("resource", resource).
			Msg("Failing resource `Singularity`")
	}

	err = c.SyncCloudEventHandler(cloudEventHandler)
	if nil != err {
		log.Fatal().
			Err(err).
			Str("singularity", env.Name).
			Strs("subjects", cloudEventHandler.Spec.Subjects).
			Msg("Failing resource `CloudEventHandler`")
	}

	err = c.SyncCloudEventOutput(cloudEventOutput)
	if nil != err {
		log.Fatal().
			Err(err).
			Msg("Failing resource `CloudEventOutput`")
	}

	err = c.SyncCloudEventValidator(cloudEventValidator)
	if nil != err {
		log.Fatal().
			Err(err).
			Msg("Failing resource `CloudEventValidator`")
	}

	if err = c.Run(stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func main() {
	err := envconfig.InitWithPrefix(&env, "EventHorizon")
	if err != nil {
		rawlog.Fatalf("Failed to set environment variables: %v", err)
	}

	zerolog.SetGlobalLevel(zerolog.Disabled)

	if env.Logging.Enabled {
		level, err := zerolog.ParseLevel(env.Logging.Level)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Setting log level")
		}

		zerolog.SetGlobalLevel(level)
	}

	if env.Logging.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	switch env.Mode {
	case "standalone":
		standalone()

	case "kubernetes":
		kubernetes()
	}
}
