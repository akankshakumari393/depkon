package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Progress",type=string,JSONPath=`.status.progress`
type Depkon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DepkonSpec   `json:"spec,omitempty"`
	Status            DepkonStatus `json:"status,omitempty"`
}

type DepkonSpec struct {
	ConfigmapRef  string   `json:"configmapRef,omitempty"`
	DeploymentRef []string `json:"deploymentRef,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DepkonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Depkon `json:"items,omitempty"`
}

type DepkonStatus struct {
	Progress string `json:"progress,omitempty"`
}
