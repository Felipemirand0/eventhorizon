package controller

import (
	"context"
	"sync"
	"time"

	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/eventhorizon"
	clientset "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	acessoscheme "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	informers "acesso.io/eventhorizon/pkg/generated/informers/externalversions/eventhorizon"
	listers "acesso.io/eventhorizon/pkg/generated/listers/eventhorizon/v1alpha2"
	. "acesso.io/eventhorizon/pkg/helpers"

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
	name         string
	client       clientset.Interface
	workqueue    workqueue.RateLimitingInterface
	eventhorizon *eventhorizon.EventHorizon
	context      context.Context
	cancel       context.CancelFunc
	mutex        *sync.RWMutex
	threadiness  int
	lister       listers.EventHorizonLister
	synced       cache.InformerSynced
}

func NewStandalone(name string) *Controller {
	ctx, cancel := context.WithCancel(context.Background())

	return &Controller{
		name:      name,
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventHorizon"),
		context:   ctx,
		cancel:    cancel,
		mutex:     &sync.RWMutex{},
	}
}

func NewKubernetes(name string, threadiness int, client clientset.Interface, informer informers.Interface) *Controller {
	utilruntime.Must(acessoscheme.AddToScheme(scheme.Scheme))

	ctx, cancel := context.WithCancel(context.Background())

	controller := &Controller{
		name:        name,
		client:      client,
		workqueue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventHorizon"),
		context:     ctx,
		cancel:      cancel,
		mutex:       &sync.RWMutex{},
		threadiness: threadiness,
		lister:      informer.V1alpha2().EventHorizons().Lister(),
		synced:      informer.V1alpha2().EventHorizons().Informer().HasSynced,
	}

	klog.Info("Setting up event handlers")

	informer.
		V1alpha2().
		EventHorizons().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
		})

	return controller
}

func (c *Controller) Run(stopCh <-chan struct{}) error {
	log.Info().
		Str("name", c.name).
		Msg("Starting instance")

	defer func() {
		<-c.context.Done()

		log.Info().
			Str("name", c.name).
			Msg("Stopping instance")

		if nil != c.eventhorizon {
			c.eventhorizon.Close()
		}

		klog.Info("Shut down workers")
		c.workqueue.ShutDown()
	}()

	defer utilruntime.HandleCrash()

	klog.Info("Wait for informer caches to sync")

	synceds := ValidSynced([]cache.InformerSynced{
		c.synced,
	})

	if ok := cache.WaitForCacheSync(stopCh, synceds...); !ok {
		return ErrWaitCacheSync
	}

	klog.Info("Start workers")

	for i := 1; i <= c.threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	go func() {
		// wait for EventHorizon
		for {
			if nil != c.context.Err() {
				return
			}

			if nil != c.eventhorizon {
				break
			}
		}

		// start it
		c.eventhorizon.Start() // this is a blocking call

		c.cancel()
	}()

	<-stopCh
	c.cancel()

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
			return nil
		}

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if nil != err {
			return err
		}

		if nil != c.eventhorizon {
			c.workqueue.Forget(key)
			klog.Errorf("Service already initiated, live reload not implemented, skipping resource '%s'.", key)
			return nil
		}

		e, err := c.client.
			EventhorizonV1alpha2().
			EventHorizons(namespace).
			Get(name, metav1.GetOptions{})

		if nil != err {
			klog.Errorf("Failed loading resource '%s', skipping", key)
			c.workqueue.Forget(key)
			return err
		}

		err = c.SyncEventHorizon(e)

		if ErrNameMismatch == err {
			klog.Errorf("Mismatch instance name (%s) with resource name (%s), skipping", c.name, key)
			c.workqueue.Forget(key)
			return nil
		}

		if nil != err {
			klog.Errorf("Unable to sync resource '%s', retrying", key)
			c.workqueue.AddRateLimited(key)
			return nil
		}

		klog.Infof("Successfully synced resource '%s'", key)
		c.workqueue.Forget(key)
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

	c.workqueue.Add(key)
}
