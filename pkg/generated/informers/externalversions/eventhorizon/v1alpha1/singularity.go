/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	eventhorizonv1alpha1 "acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha1"
	versioned "acesso.io/eventhorizon/pkg/generated/clientset/versioned"
	internalinterfaces "acesso.io/eventhorizon/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "acesso.io/eventhorizon/pkg/generated/listers/eventhorizon/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SingularityInformer provides access to a shared informer and lister for
// Singularities.
type SingularityInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SingularityLister
}

type singularityInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSingularityInformer constructs a new informer for Singularity type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSingularityInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSingularityInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSingularityInformer constructs a new informer for Singularity type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSingularityInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventhorizonV1alpha1().Singularities(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventhorizonV1alpha1().Singularities(namespace).Watch(options)
			},
		},
		&eventhorizonv1alpha1.Singularity{},
		resyncPeriod,
		indexers,
	)
}

func (f *singularityInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSingularityInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *singularityInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&eventhorizonv1alpha1.Singularity{}, f.defaultInformer)
}

func (f *singularityInformer) Lister() v1alpha1.SingularityLister {
	return v1alpha1.NewSingularityLister(f.Informer().GetIndexer())
}