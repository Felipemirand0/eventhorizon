package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	. "acesso.io/eventhorizon/pkg/errors"
	clientset "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	acessoscheme "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	eventhorizon "acesso.io/eventhorizon/pkg/generated/informers/externalversions/eventhorizon"
	listers "acesso.io/eventhorizon/pkg/generated/listers/eventhorizon/v1alpha1"
	"acesso.io/eventhorizon/pkg/handler"
	. "acesso.io/eventhorizon/pkg/helpers"
	"acesso.io/eventhorizon/pkg/output"
	"acesso.io/eventhorizon/pkg/singularity"
	"acesso.io/eventhorizon/pkg/validator"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	name        string
	client      clientset.Interface
	workqueue   workqueue.RateLimitingInterface
	singularity *singularity.Singularity
	handlers    map[string]handler.Handler
	outputs     map[string]output.Output
	validators  map[string]validator.Validator
	context     context.Context
	cancel      context.CancelFunc
	mutex       *sync.RWMutex
	threadiness int

	// kind: Singularity
	singularityLister listers.SingularityLister
	singularitySynced cache.InformerSynced

	// kind: CloudEventOutput
	cloudeventoutputLister listers.CloudEventOutputLister
	cloudeventoutputSynced cache.InformerSynced

	// kind: CloudEventHandler
	cloudeventhandlerLister listers.CloudEventHandlerLister
	cloudeventhandlerSynced cache.InformerSynced

	// kind: CloudEventValidator
	cloudeventvalidatorLister listers.CloudEventValidatorLister
	cloudeventvalidatorSynced cache.InformerSynced
}

func NewStandalone(name string) *Controller {
	ctx, cancel := context.WithCancel(context.Background())

	return &Controller{
		name:       name,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventHorizon"),
		handlers:   map[string]handler.Handler{},
		outputs:    map[string]output.Output{},
		validators: map[string]validator.Validator{},
		context:    ctx,
		cancel:     cancel,
		mutex:      &sync.RWMutex{},
	}
}

func NewKubernetes(name string, threadiness int, client clientset.Interface, eventhorizon eventhorizon.Interface) *Controller {
	utilruntime.Must(acessoscheme.AddToScheme(scheme.Scheme))

	ctx, cancel := context.WithCancel(context.Background())

	controller := &Controller{
		name:        name,
		client:      client,
		workqueue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventHorizon"),
		handlers:    map[string]handler.Handler{},
		outputs:     map[string]output.Output{},
		validators:  map[string]validator.Validator{},
		context:     ctx,
		cancel:      cancel,
		mutex:       &sync.RWMutex{},
		threadiness: threadiness,

		singularityLister: eventhorizon.V1alpha1().Singularities().Lister(),
		singularitySynced: eventhorizon.V1alpha1().Singularities().Informer().HasSynced,

		cloudeventoutputLister: eventhorizon.V1alpha1().CloudEventOutputs().Lister(),
		cloudeventoutputSynced: eventhorizon.V1alpha1().CloudEventOutputs().Informer().HasSynced,

		cloudeventhandlerLister: eventhorizon.V1alpha1().CloudEventHandlers().Lister(),
		cloudeventhandlerSynced: eventhorizon.V1alpha1().CloudEventHandlers().Informer().HasSynced,

		cloudeventvalidatorLister: eventhorizon.V1alpha1().CloudEventValidators().Lister(),
		cloudeventvalidatorSynced: eventhorizon.V1alpha1().CloudEventValidators().Informer().HasSynced,
	}

	klog.Info("Setting up event handlers")

	eventhorizon.
		V1alpha1().
		Singularities().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
		})

	eventhorizon.
		V1alpha1().
		CloudEventOutputs().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
		})

	eventhorizon.
		V1alpha1().
		CloudEventHandlers().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
		})

	eventhorizon.
		V1alpha1().
		CloudEventValidators().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
		})

	return controller
}

func (c *Controller) Run(stopCh <-chan struct{}) error {
	log.Info().
		Str("name", c.name).
		Msg("Create black hole singularity")

	klog.Info("Start controller")

	defer utilruntime.HandleCrash()

	defer func() {
		<-c.context.Done()

		time.Sleep(3 * time.Second)

		log.Info().
			Str("name", c.name).
			Msg("Evaporating black hole singularity")

		c.workqueue.ShutDown()
	}()

	klog.Info("Wait for informer caches to sync")

	synceds := ValidSynced([]cache.InformerSynced{
		c.singularitySynced,
		c.cloudeventoutputSynced,
		c.cloudeventhandlerSynced,
		c.cloudeventvalidatorSynced,
	})

	if ok := cache.WaitForCacheSync(stopCh, synceds...); !ok {
		return ErrWaitCacheSync
	}

	klog.Info("Start workers")

	for i := 1; i <= c.threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh

	c.cancel()

	klog.Info("Shut down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)

		key, ok := obj.(string)
		if !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))

			return nil
		}

		kind, namespace, name := ParseKey(key)

		var err error

		switch kind {
		case "singularity":
			e, err := c.client.
				EventhorizonV1alpha1().
				Singularities(namespace).
				Get(name, metav1.GetOptions{})
			if nil != err {
				return err
			}

			err = c.SyncSingularity(e)
			if ErrNameMismatch == err {
				c.workqueue.Forget(key)
				klog.Errorf("current singularity name '%s' mismatch resource '%s/%s', skipping", c.name, namespace, name)

				return nil
			}

		case "cloudeventoutput":
			e, err := c.client.
				EventhorizonV1alpha1().
				CloudEventOutputs().
				Get(name, metav1.GetOptions{})
			if nil != err {
				return err
			}

			err = c.SyncCloudEventOutput(e)

		case "cloudeventhandler":
			e, err := c.client.
				EventhorizonV1alpha1().
				CloudEventHandlers().
				Get(name, metav1.GetOptions{})
			if nil != err {
				return err
			}

			err = c.SyncCloudEventHandler(e)
			if ErrNoMatchingSubject == err {
				c.workqueue.Forget(key)
				klog.Errorf("no registered instance for the handler '%s', skipping", name)

				return nil
			}

		case "cloudeventvalidator":
			e, err := c.client.
				EventhorizonV1alpha1().
				CloudEventValidators().
				Get(name, metav1.GetOptions{})
			if nil != err {
				return err
			}

			err = c.SyncCloudEventValidator(e)

		default:
			c.workqueue.Forget(key)
			klog.Errorf("not a monitored resource '%s', skipping", key)

			return nil
		}

		if nil != err {
			c.workqueue.AddRateLimited(key)

			return err
		}

		c.workqueue.Forget(key)
		klog.Infof("Successfully synced '%s'", key)

		return nil
	}(obj)

	if nil != err {
		utilruntime.HandleError(err)

		return false
	}

	return true
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if nil != err {
		utilruntime.HandleError(err)

		return
	}

	switch obj.(type) {
	case *v1alpha1.Singularity:
		key = fmt.Sprintf(`%s#%s`, "singularity", key)

	case *v1alpha1.CloudEventOutput:
		key = fmt.Sprintf(`%s#%s`, "cloudeventoutput", key)

	case *v1alpha1.CloudEventHandler:
		key = fmt.Sprintf(`%s#%s`, "cloudeventhandler", key)

	case *v1alpha1.CloudEventValidator:
		key = fmt.Sprintf(`%s#%s`, "cloudeventvalidator", key)
	}

	c.workqueue.Add(key)
}
