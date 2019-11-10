package main

import (
	"context"
	"io/ioutil"
	rawlog "log"
	"os"
	"time"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"
	"acesso.io/eventhorizon/pkg/controller"
	clientset "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	acessoschema "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	informers "acesso.io/eventhorizon/pkg/generated/informers/externalversions"
	"acesso.io/eventhorizon/pkg/signals"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
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

func kubernetes(ctx context.Context) {
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

	c := controller.New(ctx, env.Name, env.Threadiness, client, informerFactory.Eventhorizon())

	informerFactory.Start(ctx.Done())

	if err = c.Run(); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func standalone(ctx context.Context) {
	acessoschema.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode

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

	var r *v1alpha2.EventHorizon

	switch obj.(type) {
	case *v1alpha2.EventHorizon:
		r = obj.(*v1alpha2.EventHorizon)

	default:
		log.Fatal().
			Msg("Configuration resource file must be of kind `EventHorizon`")
	}

	c := controller.NewStandalone(ctx, env.Name, r)

	if err = c.Run(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Error running standalone controller")
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

	exit := signals.SetupSignalHandler()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-exit
		cancel()
		<-ctx.Done()
	}()

	switch env.Mode {
	case "standalone":
		standalone(ctx)

	case "kubernetes":
		kubernetes(ctx)
	}
}
