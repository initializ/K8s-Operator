package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// InitzSecret is the Schema for the initzsecrets API
type InitzSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InitzSecretSpec   `json:"spec,omitempty"`
	Status InitzSecretStatus `json:"status,omitempty"`
}

// InitzSecretSpec defines the desired state of InitzSecret
type InitzSecretSpec struct {
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
	Workspace  string   `json:"workspace,omitempty"`
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

// +kubebuilder:object:root=true

// InitzSecretList contains a list of InitzSecret
type InitzSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InitzSecret `json:"items"`
}

// DeepCopyObject implements runtime.Object.
func (i *InitzSecretList) DeepCopyObject() runtime.Object {
	panic("unimplemented")
}

// GetObjectKind implements runtime.Object.
// Subtle: this method shadows the method (TypeMeta).GetObjectKind of InitzSecretList.TypeMeta.
func (i *InitzSecretList) GetObjectKind() schema.ObjectKind {
	panic("unimplemented")
}

// InitzSecretStatus defines the observed state of InitzSecret
type InitzSecretStatus struct {
	LastReconcileTime metav1.Time `json:"lastReconcileTime,omitempty"`
}

func (in *InitzSecret) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

// GetObjectKind returns the type of the object
func (in *InitzSecret) GetObjectKind() schema.ObjectKind {
	return &in.TypeMeta
}

// GetObjectMeta returns the object metadata
func (in *InitzSecret) GetObjectMeta() metav1.Object {
	return &in.ObjectMeta
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *InitzSecret) DeepCopyInto(out *InitzSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// Implement client.Object interface for InitzSecret
func (in *InitzSecret) GetNamespace() string {
	return in.Namespace
}

func (in *InitzSecret) SetNamespace(namespace string) {
	in.Namespace = namespace
}

func (in *InitzSecret) GetLabels() map[string]string {
	return in.Labels
}

func (in *InitzSecret) SetLabels(labels map[string]string) {
	in.Labels = labels
}

func (in *InitzSecret) DeepCopy() *InitzSecret {
	out := &InitzSecret{}
	in.DeepCopyInto(out)
	return out
}
