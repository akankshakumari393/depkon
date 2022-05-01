package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	klientset "github.com/akankshakumari393/depkon/pkg/generated/clientset/versioned"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	homeDir := os.Getenv("HOME")
	kubeconfigFile := homeDir + "/.kube/config"
	kubeconfig := flag.String("kubeconfig", kubeconfigFile, "Kubeconfig File location")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	klientsetconfig, err := klientset.NewForConfig(config)
	if err != nil {
		// handle error
		fmt.Printf("error %s, depkonclientset\n", err.Error())
	}
	_, err = klientsetconfig.Akankshakumari393V1alpha1().Depkons("default").List(context.Background(), v1.ListOptions{})
	if err != nil {
		// handle error
		fmt.Printf("error %s, depkonclientset\n", err.Error())
	}

}
