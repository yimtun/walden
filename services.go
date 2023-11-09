package main

// 让集群外部服务可以使用  cluster ip
// 让集群外部服务可感知    cluster ip
// 负载均衡 1 iptables dant 2 kong 3 lvm 4  ipvs    5 f5
import (
	"flag"
	"fmt"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"path/filepath"
	"time"

	//"k8s.io/client-go/pkg/api/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// 自定义业务

//
type ServiceEndPointLoggingController struct {
	informerFactory informers.SharedInformerFactory
	svcInformer     coreinformers.ServiceInformer
	epInformer      coreinformers.EndpointsInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *ServiceEndPointLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache. // 同时处理两个informer   ServiceInformer  EndpointsInformer
	if !cache.WaitForCacheSync(stopCh, c.svcInformer.Informer().HasSynced, c.epInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}

func NewSvcEpLoggingController(informerFactory informers.SharedInformerFactory) *ServiceEndPointLoggingController {
	svcInformer := informerFactory.Core().V1().Services()
	epInformer := informerFactory.Core().V1().Endpoints()

	c := &ServiceEndPointLoggingController{
		informerFactory: informerFactory,
		svcInformer:     svcInformer,
		epInformer:      epInformer,
	}

	// 处理 svc informer 事件  过滤掉不是 default ns 的事件
	svcInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				newSvc := obj.(*v1.Service)
				if newSvc.Namespace != "default" {
					return false
				}
				klog.Infof("filter: svc [%s/%s/%s]\n", newSvc.Namespace, "service 名字"+newSvc.Name, "service cluster ip "+newSvc.Spec.ClusterIP)
				//  自定义逻辑 将cluster ip endpoint endpoint 后的 podip 传给上层业务 将数据写入redis

				//newSvc.Spec.ClusterIP

				//
				return true
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					newSvc := obj.(*v1.Service)
					klog.Infof("controller: add svc, svc [%s/%s]\n", newSvc.Namespace, newSvc.Name)
				},

				UpdateFunc: func(oldObj, newObj interface{}) {
					newSvc := newObj.(*v1.Service)
					klog.Infof("controller: Update svc, pod [%s/%s]\n", newSvc.Namespace, newSvc.Name)
				},

				DeleteFunc: func(obj interface{}) {
					delSvc := obj.(*v1.Service)
					klog.Infof("controller: Delete svc, pod [%s/%s]\n", delSvc.Namespace, delSvc.Name)
				},
			},
		},
	)

	// 处理 endpoints  informer 事件 过滤掉不是 default ns 的事件
	epInformer.Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				ep := obj.(*v1.Endpoints)
				if ep.Namespace != "default" {
					return false
				}
				return true
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: func(cur interface{}) {
					endpoint := cur.(*v1.Endpoints)
					klog.Infof("AddEp:%v:%v", endpoint.Name, endpoint.Subsets[0].Ports[0].Port)
					//klog.Infof("AddEp:%v", endpoint.Name, &endpoint.Subsets[0].Ports, endpoint.Subsets[0].Addresses)

				},

				UpdateFunc: func(objA, objB interface{}) {
					ep1 := objA.(*v1.Endpoints)
					ep2 := objB.(*v1.Endpoints)
					klog.Infof("UpdateEp, name:%s, oldEp:%v, newEp:%v", ep1.Name, ep1.Subsets, ep2.Subsets)

				},

				DeleteFunc: func(cur interface{}) {
					endpoint := cur.(*v1.Endpoints)
					klog.Infof("DelEp [%s/%s]\n", endpoint.Namespace, endpoint.Name)
				},
			},
		},
	)

	return c
}

func main() {

	go serviceRsync()

	//  自定义业务逻辑
	//
	//1 获取service 信息   cluster ip  地址到本机  netlink 库
	//2 设置负载均衡参数，可以是iptables 规则的 目的地址替换  将访问本机ip 的服务  负载均衡方式转发给后端 pod ip

	//
	select {} // 一直阻塞 不会退出
}

func serviceRsync() {

	flag.Parse()
	logs.InitLogs()
	defer logs.FlushLogs()
	var kubeconfigTemp *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath := filepath.Join(home, ".kube", "config")
		kubeconfigTemp = flag.String("kubeconfig1", kubeconfigPath, "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigTemp = flag.String("kubeconfig1", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigTemp)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller := NewSvcEpLoggingController(factory)

	stop := make(chan struct{})
	defer close(stop)

	err = controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}

	// 此处一直阻塞不退出
	select {}
}
