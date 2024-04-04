/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InitzSecret is the Schema for the initzsecrets API
type InitzSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec InitzSecretSpec `json:"spec,omitempty"`
}

// InitzSecretSpec defines the desired state of InitzSecret
type InitzSecretSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	HostAPI                string                 `json:"hostAPI,omitempty"`
	ResyncInterval         int64                  `json:"resyncInterval,omitempty"`
	Authentication         Authentication         `json:"authentication,omitempty"`
	ManagedSecretReference ManagedSecretReference `json:"managedSecretReference,omitempty"`
}

// Authentication defines the authentication details
type Authentication struct {
	ServiceToken ServiceToken `json:"serviceToken,omitempty"`
}

// ServiceToken defines the service token details
type ServiceToken struct {
	ServiceTokenSecretReference SecretReference `json:"serviceTokenSecretReference,omitempty"`
	SecretsScope                SecretsScope    `json:"secretsScope,omitempty"`
}

// SecretsScope defines the scope for fetching secrets
type SecretsScope struct {
	OrganisationID  string   `json:"organisationID,omitempty"`
	EnvSlug    string   `json:"envSlug,omitempty"`
	SecretVars []string `json:"secretVars,omitempty"`
}

// SecretReference defines the reference to a Kubernetes secret
type SecretReference struct {
	Servicetoken string `json:"servicetoken,omitempty"`
}

// ManagedSecretReference defines the reference to the managed Kubernetes secret
type ManagedSecretReference struct {
	SecretName      string `json:"secretName,omitempty"`
	SecretNamespace string `json:"secretNamespace,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

//+kubebuilder:object:root=true

// InitzSecretList contains a list of InitzSecret
type InitzSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InitzSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InitzSecret{}, &InitzSecretList{})
}

func (in *InitzSecret) DeepCopyObject() runtime.Object {
	out := &InitzSecret{}
	in.DeepCopyInto(out)
	return out
}
