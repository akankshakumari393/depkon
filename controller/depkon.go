package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"

	"github.com/akankshakumari393/depkon/pkg/apis/akankshakumari393.dev/v1alpha1"
	customscheme "github.com/akankshakumari393/depkon/pkg/generated/clientset/versioned/scheme"

	klientset "github.com/akankshakumari393/depkon/pkg/generated/clientset/versioned"
	klientinformer "github.com/akankshakumari393/depkon/pkg/generated/informers/externalversions/akankshakumari393.dev/v1alpha1"
	klientlister "github.com/akankshakumari393/depkon/pkg/generated/listers/akankshakumari393.dev/v1alpha1"
)

type controller struct {
	// clientset for custom resource depkon
	klientset klientset.Interface
	// clientset for native resource
	clientset kubernetes.Interface
	kLister   klientlister.DepkonLister
	// depkon has synced
	depkonSynced cache.InformerSynced
	queue        workqueue.RateLimitingInterface
	// record events
	recorder record.EventRecorder
}

func NewController(clientset kubernetes.Interface, klientset klientset.Interface, depkonInformer klientinformer.DepkonInformer) *controller {
	runtime.Must(customscheme.AddToScheme(scheme.Scheme))
	eveBroadcaster := record.NewBroadcaster()
	eveBroadcaster.StartStructuredLogging(1)
	eveBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: clientset.CoreV1().Events(""),
	})
	recorder := eveBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "Depkon"})
	controller := &controller{
		clientset:    clientset,
		klientset:    klientset,
		kLister:      depkonInformer.Lister(),
		depkonSynced: depkonInformer.Informer().HasSynced,
		queue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "env-sync"),
		recorder:     recorder,
	}

	depkonInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    controller.depkonAdded,
			UpdateFunc: controller.depkonUpdated,
			DeleteFunc: controller.depkonDeleted,
		},
	)
	return controller
}

func (c *controller) Run(numberOfWorkers int, stopCh <-chan struct{}) error {
	log.Println("Running controller")
	// Informer maintains a cache that needs to be synced for the first time this brings configMap from default namespace
	if !cache.WaitForCacheSync(stopCh, c.depkonSynced) {
		log.Println("error Cache not synced")
	}

	log.Println("Starting workers")
	for i := 0; i < numberOfWorkers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}
	log.Println("Started workers")
	<-stopCh
	log.Println("Shutting down workers")
	return nil
}

// worker function operate on the queue and process each item
func (c *controller) worker() {
	log.Println("Worker called")
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("\nError getting key from Item %s", err.Error())
		return false
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("\nError split key from Item %s", err.Error())
		return false
	}
	depkonResource, err := c.kLister.Depkons(ns).Get(name)
	if err != nil {
		c.recorder.Event(depkonResource, v1.EventTypeWarning, "DepkonFailed", "Failed to get depkon resource from lister")
		return false
	}
	c.recorder.Event(depkonResource, v1.EventTypeNormal, "DepkonCreated", fmt.Sprintf("Sync Depkon resource %s from namespace %s", depkonResource.Name, depkonResource.Namespace))
	c.updateStatus("Syncing", depkonResource)
	err = c.syncDepkon(depkonResource)
	if err != nil {
		c.recorder.Event(depkonResource, v1.EventTypeWarning, "DepkonSyncFailed", fmt.Sprintf("Failed syncing Depkon resource %s from namespace %s", depkonResource.Name, depkonResource.Namespace))
		return false
	}
	c.updateStatus("Synced", depkonResource)
	return true
}

func (c *controller) updateStatus(progress string, depkonResource *v1alpha1.Depkon) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Depkon before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := c.klientset.Akankshakumari393V1alpha1().Depkons(depkonResource.Namespace).Get(context.Background(), depkonResource.Name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		result.Status.Progress = progress
		_, updateErr := c.klientset.Akankshakumari393V1alpha1().Depkons(depkonResource.Namespace).UpdateStatus(context.Background(), result, metav1.UpdateOptions{})

		if updateErr == nil {
			c.recorder.Event(depkonResource, v1.EventTypeNormal, "StatusUpdated", fmt.Sprintf("Status set to %s", progress))
		}
		return updateErr
	})
	if retryErr != nil {
		c.recorder.Event(depkonResource, v1.EventTypeWarning, "FailedUpdatingStatus", fmt.Sprintf("Failed to Update status %s", retryErr.Error()))
		return retryErr
	}
	return nil
}

func (c *controller) syncDepkon(depkonResource *v1alpha1.Depkon) error {
	ctx := context.Background()

	// update the deployments, Add the configMapRef
	for _, deploymentName := range depkonResource.Spec.DeploymentRef {
		c.recorder.Event(depkonResource, v1.EventTypeNormal, "UpdateDeployment", fmt.Sprintf("%s updating deployment with configmap %s", deploymentName, depkonResource.Spec.ConfigmapRef))
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := c.clientset.AppsV1().Deployments(depkonResource.Namespace).Get(ctx, deploymentName, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
			addCmRef := true
			for _, envRef := range result.Spec.Template.Spec.Containers[0].EnvFrom {
				if envRef.ConfigMapRef.LocalObjectReference.Name == depkonResource.Spec.ConfigmapRef {
					addCmRef = false
				}
			}
			if !addCmRef {
				c.recorder.Event(depkonResource, v1.EventTypeNormal, "DeploymentAlreadyUpdated", fmt.Sprintf("%s Deployment with configmap already updated %s", deploymentName, depkonResource.Spec.ConfigmapRef))
				return nil
			}

			EnvFrom := v1.EnvFromSource{
				ConfigMapRef: &v1.ConfigMapEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: depkonResource.Spec.ConfigmapRef,
					},
				},
			}
			result.Spec.Template.Spec.Containers[0].EnvFrom = append(result.Spec.Template.Spec.Containers[0].EnvFrom, EnvFrom)
			_, updateErr := c.clientset.AppsV1().Deployments(depkonResource.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
			if updateErr == nil {
				c.recorder.Event(depkonResource, v1.EventTypeNormal, "DeploymentUpdated", fmt.Sprintf("%s Deployment with configmap updated %s", deploymentName, depkonResource.Spec.ConfigmapRef))
			}
			return updateErr
		})
		if retryErr != nil {
			c.recorder.Event(depkonResource, v1.EventTypeWarning, "FailedUpdatingDeployment", fmt.Sprintf("%s deployment failed to be updated %s", deploymentName, retryErr.Error()))
			return retryErr
		}
	}
	return nil
}

func (c *controller) depkonAdded(obj interface{}) {
	log.Println("Depkon Added")
	c.queue.Add(obj)
}
func (c *controller) depkonDeleted(obj interface{}) {
	log.Println("Depkon Deleted")
	// c.queue.Add(obj)
}
func (c *controller) depkonUpdated(oldobj, newobj interface{}) {
	log.Println("Depkon Updated")
	//c.queue.Add(newobj)
}
