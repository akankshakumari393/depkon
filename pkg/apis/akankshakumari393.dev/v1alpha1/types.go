package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Depkon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DepkonSpec `json:"spec,omitempty"`
}

type DepkonSpec struct {
	Config []NamespaceConfig `json:"config,omitempty"`
}

type NamespaceConfig struct {
	ConfigmapRef  string   `json:"configmapRef,omitempty"`
	Namespace     string   `json:"namespace,omitempty"`
	DeploymentRef []string `json:"deploymentRef,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DepkonList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []Depkon `json:"items,omitempty"`
}
