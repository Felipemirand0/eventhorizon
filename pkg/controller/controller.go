package controller

import (
	"context"
	"time"

	. "acesso.io/eventhorizon/pkg/errors"
	"acesso.io/eventhorizon/pkg/eventhorizon"
	clientset "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	acessoscheme "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	informers "acesso.io/eventhorizon/pkg/generated/informers/externalversions/eventhorizon"
	listers "acesso.io/eventhorizon/pkg/generated/listers/eventhorizon/v1alpha2"
	. "acesso.io/eventhorizon/pkg/helpers"

	"github.com/rs/zerolog"
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
	threadiness  int
	lister       listers.EventHorizonLister
	synced       cache.InformerSynced
}

func New(ctx context.Context, name string, threadiness int, client clientset.Interface, informer informers.Interface) *Controller {
	utilruntime.Must(acessoscheme.AddToScheme(scheme.Scheme))

	c := &Controller{
		name:        name,
		client:      client,
		workqueue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventHorizon"),
		context:     ctx,
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
			AddFunc: c.enqueue,
		})

	return c
}

func (c *Controller) Run() error {
	log.Info().
		Str("name", c.name).
		Msg("Starting controller")

	klog.Info("Wait for informer caches to sync")

	synceds := ValidSynced([]cache.InformerSynced{
		c.synced,
	})

	if ok := cache.WaitForCacheSync(c.context.Done(), synceds...); !ok {
		return ErrWaitCacheSync
	}

	klog.Info("Start workers")

	for i := 1; i <= c.threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, c.context.Done())
	}

	<-c.context.Done()

	log.Info().
		Str("name", c.name).
		Msg("Stopping controller")

	if nil != c.eventhorizon {
		c.eventhorizon.Close()
	}

	klog.Info("Stop workers")
	c.workqueue.ShutDown()

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
		if nil != c.eventhorizon {
			c.workqueue.ShutDown()
		}
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

		r, err := c.client.
			EventhorizonV1alpha2().
			EventHorizons(namespace).
			Get(name, metav1.GetOptions{})

		if nil != err {
			klog.Errorf("Failed loading resource '%s', skipping", key)
			c.workqueue.Forget(key)
			return err
		}

		e, err := c.syncEventHorizon(r)

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

		c.eventhorizon = e

		go func() {
			log.Info().
				Msg("Starting EventHorizon")

			err := c.eventhorizon.Start()

			c.eventhorizon.Close()

			var level zerolog.Level = zerolog.InfoLevel

			if nil != err {
				level = zerolog.ErrorLevel
				c.workqueue.AddRateLimited(key)
			}

			log.WithLevel(level).
				Err(err).
				Msg("Stopping EventHorizon")
		}()

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
