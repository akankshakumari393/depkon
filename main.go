package main

import (
	"flag"
	"fmt"
	"os"

	// klientset "github.com/akankshakumari393/depkon/pkg/client/clientset/versioned"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	homeDir := os.Getenv("HOME")
	kubeconfigFile := homeDir + "/.kube/config"
	kubeconfig := flag.String("kubeconfig", kubeconfigFile, "Kubeconfig File location")
	_, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		_, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	// depkonclientset, err := klientset.NewForConfig(config)
	// if err != nil {
	// 	// handle error
	// 	fmt.Printf("error %s, depkonclientset\n", err.Error())
	// }

}
