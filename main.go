package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	kcontroller "github.com/akankshakumari393/depkon/controller"
	klientset "github.com/akankshakumari393/depkon/pkg/generated/clientset/versioned"
	klientfactory "github.com/akankshakumari393/depkon/pkg/generated/informers/externalversions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		log.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	klientset, err := klientset.NewForConfig(config)
	if err != nil {
		// handle error
		log.Printf("error %s, depkonclientset\n", err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// handle error
		log.Printf("error %s, kubernetes clientset\n", err.Error())
	}

	_, err = klientset.Akankshakumari393V1alpha1().Depkons("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// handle error
		log.Printf("error %s, depkonclientset\n", err.Error())
	}

	//If we want to create informer for specific resources
	// labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
	// 	//opts.LabelSelector = "app=nats-box"
	// })

	// By default NewSharedInformerFactory creates informerfactory for all Namespaces
	// Use NewSharedInformerFactoryWithOptions for creating informer instance in specific namespace
	infofactory := klientfactory.NewSharedInformerFactory(klientset, 10*time.Minute)

	ch := make(chan struct{})

	depkonController := kcontroller.NewController(clientset, klientset, infofactory.Akankshakumari393().V1alpha1().Depkons())

	infofactory.Start(ch)

	if err := depkonController.Run(3, ch); err != nil {
		log.Printf("error running contoller %s\n", err.Error())
	}

}
