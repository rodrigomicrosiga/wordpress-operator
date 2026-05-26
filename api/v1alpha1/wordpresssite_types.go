package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec define a configuração do banco de dados MySQL
type DatabaseSpec struct {
	// +kubebuilder:default:="mysql:8.0"
	// +optional
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	User string `json:"user"`

	// +kubebuilder:default:="5Gi"
	// +optional
	StorageSize string `json:"storageSize,omitempty"`
}

// WordpressSpec define a configuração da aplicação WordPress
type WordpressSpec struct {
	// +kubebuilder:default:="wordpress:6.5-apache"
	// +optional
	Image string `json:"image,omitempty"`

	// +kubebuilder:default:=1
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:default:="2Gi"
	// +optional
	StorageSize string `json:"storageSize,omitempty"`
}

// WordpressSiteSpec define o estado desejado pelo usuário (Spec)
type WordpressSiteSpec struct {
	// +kubebuilder:validation:Required
	Domain string `json:"domain"`

	// +kubebuilder:validation:Required
	Wordpress WordpressSpec `json:"wordpress"`

	// +kubebuilder:validation:Required
	Database DatabaseSpec `json:"database"`

	// +optional
	IngressClassName string `json:"ingressClassName,omitempty"`
}

// WordpressSiteStatus define o estado observado pelo Operator (Status)
type WordpressSiteStatus struct {
	// +optional
	Phase string `json:"phase,omitempty"`

	// +optional
	URL string `json:"url,omitempty"`

	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Domain",type="string",JSONPath=".spec.domain"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// WordpressSite is the Schema for the wordpresssites API
type WordpressSite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WordpressSiteSpec   `json:"spec,omitempty"`
	Status WordpressSiteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WordpressSiteList contains a list of WordpressSite
type WordpressSiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WordpressSite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WordpressSite{}, &WordpressSiteList{})
}
