// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventHandler) DeepCopyInto(out *CloudEventHandler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventHandler.
func (in *CloudEventHandler) DeepCopy() *CloudEventHandler {
	if in == nil {
		return nil
	}
	out := new(CloudEventHandler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventHandler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventHandlerList) DeepCopyInto(out *CloudEventHandlerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CloudEventHandler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventHandlerList.
func (in *CloudEventHandlerList) DeepCopy() *CloudEventHandlerList {
	if in == nil {
		return nil
	}
	out := new(CloudEventHandlerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventHandlerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventHandlerSpec) DeepCopyInto(out *CloudEventHandlerSpec) {
	*out = *in
	if in.Subjects != nil {
		in, out := &in.Subjects, &out.Subjects
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventHandlerSpec.
func (in *CloudEventHandlerSpec) DeepCopy() *CloudEventHandlerSpec {
	if in == nil {
		return nil
	}
	out := new(CloudEventHandlerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventOutput) DeepCopyInto(out *CloudEventOutput) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventOutput.
func (in *CloudEventOutput) DeepCopy() *CloudEventOutput {
	if in == nil {
		return nil
	}
	out := new(CloudEventOutput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventOutput) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventOutputFluentd) DeepCopyInto(out *CloudEventOutputFluentd) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventOutputFluentd.
func (in *CloudEventOutputFluentd) DeepCopy() *CloudEventOutputFluentd {
	if in == nil {
		return nil
	}
	out := new(CloudEventOutputFluentd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventOutputList) DeepCopyInto(out *CloudEventOutputList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CloudEventOutput, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventOutputList.
func (in *CloudEventOutputList) DeepCopy() *CloudEventOutputList {
	if in == nil {
		return nil
	}
	out := new(CloudEventOutputList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventOutputList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventOutputSpec) DeepCopyInto(out *CloudEventOutputSpec) {
	*out = *in
	out.Fluentd = in.Fluentd
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventOutputSpec.
func (in *CloudEventOutputSpec) DeepCopy() *CloudEventOutputSpec {
	if in == nil {
		return nil
	}
	out := new(CloudEventOutputSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventValidator) DeepCopyInto(out *CloudEventValidator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventValidator.
func (in *CloudEventValidator) DeepCopy() *CloudEventValidator {
	if in == nil {
		return nil
	}
	out := new(CloudEventValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventValidator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventValidatorList) DeepCopyInto(out *CloudEventValidatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CloudEventValidator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventValidatorList.
func (in *CloudEventValidatorList) DeepCopy() *CloudEventValidatorList {
	if in == nil {
		return nil
	}
	out := new(CloudEventValidatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudEventValidatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudEventValidatorSpec) DeepCopyInto(out *CloudEventValidatorSpec) {
	*out = *in
	if in.Handlers != nil {
		in, out := &in.Handlers, &out.Handlers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedTypes != nil {
		in, out := &in.AllowedTypes, &out.AllowedTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedSources != nil {
		in, out := &in.AllowedSources, &out.AllowedSources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudEventValidatorSpec.
func (in *CloudEventValidatorSpec) DeepCopy() *CloudEventValidatorSpec {
	if in == nil {
		return nil
	}
	out := new(CloudEventValidatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPTransport) DeepCopyInto(out *HTTPTransport) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPTransport.
func (in *HTTPTransport) DeepCopy() *HTTPTransport {
	if in == nil {
		return nil
	}
	out := new(HTTPTransport)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NATSTransport) DeepCopyInto(out *NATSTransport) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NATSTransport.
func (in *NATSTransport) DeepCopy() *NATSTransport {
	if in == nil {
		return nil
	}
	out := new(NATSTransport)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Singularity) DeepCopyInto(out *Singularity) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Singularity.
func (in *Singularity) DeepCopy() *Singularity {
	if in == nil {
		return nil
	}
	out := new(Singularity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Singularity) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SingularityList) DeepCopyInto(out *SingularityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Singularity, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SingularityList.
func (in *SingularityList) DeepCopy() *SingularityList {
	if in == nil {
		return nil
	}
	out := new(SingularityList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SingularityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SingularitySpec) DeepCopyInto(out *SingularitySpec) {
	*out = *in
	in.Transport.DeepCopyInto(&out.Transport)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SingularitySpec.
func (in *SingularitySpec) DeepCopy() *SingularitySpec {
	if in == nil {
		return nil
	}
	out := new(SingularitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Transport) DeepCopyInto(out *Transport) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(HTTPTransport)
		**out = **in
	}
	if in.NATS != nil {
		in, out := &in.NATS, &out.NATS
		*out = new(NATSTransport)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Transport.
func (in *Transport) DeepCopy() *Transport {
	if in == nil {
		return nil
	}
	out := new(Transport)
	in.DeepCopyInto(out)
	return out
}