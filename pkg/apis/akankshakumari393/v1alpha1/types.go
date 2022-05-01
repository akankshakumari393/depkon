package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
