package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Depkon struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec DepkonSpec
}

type DepkonSpec struct {
	config []NamespaceConfig
}

type NamespaceConfig struct {
	configmap  string
	namespace  string
	deployment []string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DepkonList struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Items []Depkon
}
