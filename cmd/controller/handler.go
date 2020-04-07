package main

import (
	"time"

	testresourcev1beta1 "insujang.github.io/kubernetes-test-controller/lib/testresource/v1beta1"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

// Run first sync client-go cache by calling cache.WaitforCacheSync,
// then block the main thread forever.
// Event handler will run in another Goroutine,
// generated in c.NewController() function.
func (c *Controller) Run() {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Infoln("Waiting cache to be synced.")
	// Handle timeout for syncing.
	timeout := time.NewTimer(time.Second * 30)
	timeoutCh := make(chan struct{})
	go func() {
		<-timeout.C
		timeoutCh <- struct{}{}
	}()
	if ok := cache.WaitForCacheSync(timeoutCh, c.informer.HasSynced); !ok {
		klog.Fatalln("Timeout expired during waiting for caches to sync.")
	}

	klog.Infoln("Starting custom controller.")
	select {}
}

const statusMessage = "HANDLED"

func (c *Controller) objectAddedCallback(object interface{}) {
	klog.Infof("Added: %v", object)
	resource := object.(*testresourcev1beta1.TestResource)

	// If the object is in the desired state, end callback.
	if resource.Status == statusMessage {
		return
	}

	// If the object is not handled yet, handle it by modifying its Status.
	copy := resource.DeepCopy()
	copy.Status = statusMessage
	_, err := c.testresourceclientset.InsujangV1beta1().TestResources(corev1.NamespaceDefault).Update(copy)
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	c.recorder.Event(copy, corev1.EventTypeNormal, "ObjectHandled", "Object is handled by custom controller.")
}
