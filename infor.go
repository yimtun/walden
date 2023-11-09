package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	klog "k8s.io/klog/v2"
	//"k8s.io/client-go/pkg/api/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/logs"
)

// PodLoggingController logs the name and namespace of pods that are added,
// deleted, or updated
type PodLoggingController struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

type PodLoggingController1 struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

type PodLoggingController2 struct {
	informerFactory informers.SharedInformerFactory
	serviceInformer coreinformers.ServiceInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *PodLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *PodLoggingController) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("POD CREATED: %s/%s %s", pod.Namespace, pod.Name, "创建pod时触发")
}

func (c *PodLoggingController2) podAdd(obj interface{}) {
	pod := obj.(*v1.Service)
	klog.Infof("POD CREATED: %s/%s %s", pod.Namespace, pod.Name, "创建service时触发")
}

func (c *PodLoggingController) podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	klog.Infof(
		"POD UPDATED. %s/%s %s %s",
		oldPod.Namespace, oldPod.Name, newPod.Status.Phase, "update pod时触发",
	)
}

func (c *PodLoggingController) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("POD DELETED: %s/%s %s", pod.Namespace, pod.Name, "删除pod时触发")
}

func (c *PodLoggingController1) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *PodLoggingController1) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("POD CREATED: %s/%s %s", pod.Namespace, pod.Name, "创建pod时触发")
}

func (c *PodLoggingController1) podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	klog.Infof(
		"POD UPDATED. %s/%s %s %s",
		oldPod.Namespace, oldPod.Name, newPod.Status.Phase, "update pod时触发",
	)
}

func (c *PodLoggingController1) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("POD DELETED: %s/%s %s", pod.Namespace, pod.Name, "删除pod时触发")
}

func (c *PodLoggingController2) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.serviceInformer.Informer().HasSynced) {

		return fmt.Errorf("Failed to sync")
	}
	return nil
}

// NewPodLoggingController creates a PodLoggingController
func NewPodLoggingController(informerFactory informers.SharedInformerFactory) *PodLoggingController {
	podInformer := informerFactory.Core().V1().Pods()

	c := &PodLoggingController{
		informerFactory: informerFactory,
		podInformer:     podInformer,
	}
	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}

func NewPodLoggingController1(informerFactory informers.SharedInformerFactory) *PodLoggingController1 {
	fmt.Println("NewPodLoggingController1")
	podInformer := informerFactory.Core().V1().Pods()

	c := &PodLoggingController1{
		informerFactory: informerFactory,
		podInformer:     podInformer,
	}
	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}

func NewPodLoggingController2(informerFactory informers.SharedInformerFactory) *PodLoggingController2 {
	fmt.Println("NewPodLoggingController1")
	serviceInformer := informerFactory.Core().V1().Services()

	c := &PodLoggingController2{
		informerFactory: informerFactory,
		serviceInformer: serviceInformer,
	}
	serviceInformer.Informer().AddEventHandler(

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			//UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			//DeleteFunc: c.podDelete,
		},
	)
	return c
}

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
}

func main() {
	flag.Parse()
	logs.InitLogs()
	defer logs.FlushLogs()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	// 使用一个clientset 创建了多个informmer
	go podLog(clientset)
	go podLog1(clientset)

	fmt.Println("执行到这里")

	select {} //阻塞当前goroutine
	//time.Sleep(time.Second * 30)
}

func podLog(clientset *kubernetes.Clientset) {
	fmt.Println("podLog")
	//factory := informers.NewSharedInformerFactory(clientset, time.Second*20) // 主动同步周期
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)

	controller := NewPodLoggingController(factory)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}
	select {}

}

func podLog1(clientset *kubernetes.Clientset) {
	fmt.Println("podLog1")
	//factory := informers.NewSharedInformerFactory(clientset, time.Second*10) // 主动同步周期
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)

	controller := NewPodLoggingController1(factory)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}
	select {}

}

func podLog2(clientset *kubernetes.Clientset) {
	fmt.Println("podLog1")
	//factory := informers.NewSharedInformerFactory(clientset, time.Second*10) // 主动同步周期
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)

	controller := NewPodLoggingController2(factory)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}
	select {}

}
