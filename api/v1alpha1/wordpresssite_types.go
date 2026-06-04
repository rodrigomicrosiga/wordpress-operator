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

// TLSSpec configura HTTPS no Ingress usando o cert-manager.
type TLSSpec struct {
	// Enabled liga a emissão do certificado e o bloco TLS no Ingress.
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// IssuerName é o nome do (Cluster)Issuer do cert-manager que emite o cert.
	// +optional
	IssuerName string `json:"issuerName,omitempty"`

	// IssuerKind escolhe entre "ClusterIssuer" (padrão) e "Issuer".
	// +kubebuilder:validation:Enum=ClusterIssuer;Issuer
	// +kubebuilder:default:="ClusterIssuer"
	// +optional
	IssuerKind string `json:"issuerKind,omitempty"`
}

// WordpressSiteSpec defines the desired state of WordpressSite
type WordpressSiteSpec struct {
	Domain string `json:"domain"`

	// +optional
	IngressClassName string `json:"ingressClassName,omitempty"`

	// IngressAnnotations adiciona annotations livres no Ingress.
	// +optional
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`

	// TLS configura HTTPS via cert-manager.
	// +optional
	TLS *TLSSpec `json:"tls,omitempty"`

	Wordpress WordpressSpec `json:"wordpress"`
	Database  DatabaseSpec  `json:"database"`
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
