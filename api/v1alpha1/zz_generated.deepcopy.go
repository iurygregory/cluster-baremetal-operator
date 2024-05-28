//go:build !ignore_autogenerated

/*

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnabledFeatures) DeepCopyInto(out *EnabledFeatures) {
	*out = *in
	if in.ProvisioningNetwork != nil {
		in, out := &in.ProvisioningNetwork, &out.ProvisioningNetwork
		*out = make(map[ProvisioningNetwork]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnabledFeatures.
func (in *EnabledFeatures) DeepCopy() *EnabledFeatures {
	if in == nil {
		return nil
	}
	out := new(EnabledFeatures)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreProvisioningOSDownloadURLs) DeepCopyInto(out *PreProvisioningOSDownloadURLs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreProvisioningOSDownloadURLs.
func (in *PreProvisioningOSDownloadURLs) DeepCopy() *PreProvisioningOSDownloadURLs {
	if in == nil {
		return nil
	}
	out := new(PreProvisioningOSDownloadURLs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provisioning) DeepCopyInto(out *Provisioning) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provisioning.
func (in *Provisioning) DeepCopy() *Provisioning {
	if in == nil {
		return nil
	}
	out := new(Provisioning)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Provisioning) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisioningList) DeepCopyInto(out *ProvisioningList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Provisioning, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisioningList.
func (in *ProvisioningList) DeepCopy() *ProvisioningList {
	if in == nil {
		return nil
	}
	out := new(ProvisioningList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ProvisioningList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisioningSpec) DeepCopyInto(out *ProvisioningSpec) {
	*out = *in
	if in.ProvisioningMacAddresses != nil {
		in, out := &in.ProvisioningMacAddresses, &out.ProvisioningMacAddresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.PreProvisioningOSDownloadURLs = in.PreProvisioningOSDownloadURLs
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisioningSpec.
func (in *ProvisioningSpec) DeepCopy() *ProvisioningSpec {
	if in == nil {
		return nil
	}
	out := new(ProvisioningSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisioningStatus) DeepCopyInto(out *ProvisioningStatus) {
	*out = *in
	in.OperatorStatus.DeepCopyInto(&out.OperatorStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisioningStatus.
func (in *ProvisioningStatus) DeepCopy() *ProvisioningStatus {
	if in == nil {
		return nil
	}
	out := new(ProvisioningStatus)
	in.DeepCopyInto(out)
	return out
}
