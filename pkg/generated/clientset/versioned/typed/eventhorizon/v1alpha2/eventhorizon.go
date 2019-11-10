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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	"time"

	v1alpha2 "acesso.io/eventhorizon/pkg/apis/eventhorizon/v1alpha2"
	scheme "acesso.io/eventhorizon/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EventHorizonsGetter has a method to return a EventHorizonInterface.
// A group's client should implement this interface.
type EventHorizonsGetter interface {
	EventHorizons(namespace string) EventHorizonInterface
}

// EventHorizonInterface has methods to work with EventHorizon resources.
type EventHorizonInterface interface {
	Create(*v1alpha2.EventHorizon) (*v1alpha2.EventHorizon, error)
	Update(*v1alpha2.EventHorizon) (*v1alpha2.EventHorizon, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha2.EventHorizon, error)
	List(opts v1.ListOptions) (*v1alpha2.EventHorizonList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.EventHorizon, err error)
	EventHorizonExpansion
}

// eventHorizons implements EventHorizonInterface
type eventHorizons struct {
	client rest.Interface
	ns     string
}

// newEventHorizons returns a EventHorizons
func newEventHorizons(c *EventhorizonV1alpha2Client, namespace string) *eventHorizons {
	return &eventHorizons{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the eventHorizon, and returns the corresponding eventHorizon object, and an error if there is any.
func (c *eventHorizons) Get(name string, options v1.GetOptions) (result *v1alpha2.EventHorizon, err error) {
	result = &v1alpha2.EventHorizon{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventhorizons").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EventHorizons that match those selectors.
func (c *eventHorizons) List(opts v1.ListOptions) (result *v1alpha2.EventHorizonList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.EventHorizonList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventhorizons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested eventHorizons.
func (c *eventHorizons) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("eventhorizons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a eventHorizon and creates it.  Returns the server's representation of the eventHorizon, and an error, if there is any.
func (c *eventHorizons) Create(eventHorizon *v1alpha2.EventHorizon) (result *v1alpha2.EventHorizon, err error) {
	result = &v1alpha2.EventHorizon{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("eventhorizons").
		Body(eventHorizon).
		Do().
		Into(result)
	return
}

// Update takes the representation of a eventHorizon and updates it. Returns the server's representation of the eventHorizon, and an error, if there is any.
func (c *eventHorizons) Update(eventHorizon *v1alpha2.EventHorizon) (result *v1alpha2.EventHorizon, err error) {
	result = &v1alpha2.EventHorizon{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventhorizons").
		Name(eventHorizon.Name).
		Body(eventHorizon).
		Do().
		Into(result)
	return
}

// Delete takes name of the eventHorizon and deletes it. Returns an error if one occurs.
func (c *eventHorizons) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventhorizons").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *eventHorizons) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventhorizons").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched eventHorizon.
func (c *eventHorizons) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.EventHorizon, err error) {
	result = &v1alpha2.EventHorizon{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("eventhorizons").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}