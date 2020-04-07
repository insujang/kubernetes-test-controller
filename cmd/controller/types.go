package main

import (
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	testresourceclienteset "insujang.github.io/kubernetes-test-controller/lib/testresource/generated/clientset/versioned"
	testresourecescheme "insujang.github.io/kubernetes-test-controller/lib/testresource/generated/clientset/versioned/scheme"
	testresourceinformers "insujang.github.io/kubernetes-test-controller/lib/testresource/generated/informers/externalversions"
	testresourcelisters "insujang.github.io/kubernetes-test-controller/lib/testresource/generated/listers/testresource/v1beta1"
	testresourcev1beta1 "insujang.github.io/kubernetes-test-controller/lib/testresource/v1beta1"
)

type Controller struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	testresourceclientset  testresourceclienteset.Interface
	informer               cache.SharedIndexInformer
	lister                 testresourcelisters.TestResourceLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

func NewController() *Controller {
	kubeconfig := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatalf(err.Error())
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(config)
	testClient := testresourceclienteset.NewForConfigOrDie(config)

	informerFactory := testresourceinformers.NewSharedInformerFactory(testClient, time.Minute*1)
	informer := informerFactory.Insujang().V1beta1().TestResources()

	utilruntime.Must(testresourcev1beta1.AddToScheme(testresourecescheme.Scheme))
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(testresourecescheme.Scheme, corev1.EventSource{Component: "testresource-controller"})

	workqueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	c := &Controller{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		testresourceclientset:  testClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              workqueue,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.objectAddedCallback,
		UpdateFunc: func(oldObject, newObject interface{}) {
			klog.Infof("Updated: %v", newObject)
		},
		DeleteFunc: func(object interface{}) {
			klog.Infof("Deleted: %v", object)
		},
	})
	informerFactory.Start(wait.NeverStop)

	return c
}
