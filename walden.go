package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/tools/cache"
	"log"
	"path/filepath"
	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 创建一个 client config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err)
	}

	// 初始化 client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
		panic(err)
		return
	}

	//创建 一个channel 类型是  struct{}
	stopper := make(chan struct{})
	defer close(stopper)

	//podInformer := informerFactory.Core().V1().Pods()

	// create  SharedInformerFactory

	factory := informers.NewSharedInformerFactory(clientset, 0) // 返回一个 informers.SharedInformerFactory 是一个interface
	// 其中继承了 internalinterfaces.SharedInformerFactory 这个接口

	// create  PodInformer
	podInformer := factory.Core().V1().Pods()

	//servicesInformer:=factory.Core().V1().Services().Informer()

	informer := podInformer.Informer()

	defer runtime.HandleCrash()

	// 启动 informer，list & watch

	//go factory.Start(stopper) // 开启一个goroutine  调用了internalinterfaces.SharedInformerFactory 接口里的 Start 方法

	//for {

	// 从 apiserver 同步资源，即 list
	//time.Sleep(time.Second * 30) // for 循环内部  执行时间间隔

	// 等待填数据
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) { // 传参数  channel
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	//

	// 创建 lister
	nodeLister := podInformer.Lister()

	// 从 lister 中获取所有 items
	//nodeList, err := nodeLister.List(labels.Everything())

	// watch 新创建的pod 并输出该pod 所在的node
	nodeList, err := nodeLister.Pods("default").List(labels.Everything())
	if err != nil {
		fmt.Println(err)
		return

	}

	//

	for _, node := range nodeList {
		fmt.Println("pod 数量:", len(nodeList))
		fmt.Println(node.Name)
		//fmt.Println(node.Status.ContainerStatuses)
		fmt.Println(node.Status.HostIP)
		fmt.Println(node.Status.PodIP)
		fmt.Println(node.Status.ContainerStatuses)
		if len(node.Status.ContainerStatuses) > 0 {
			//fmt.Println("打印Image:",node.Status.ContainerStatuses[0].Image)
			fmt.Println(node.Name)
			fmt.Println("打印容器ID:", node.Status.ContainerStatuses[0].ContainerID)
			fmt.Println("打印NodeIp:", node.Status.HostIP)

		}
	}
	//}

	fmt.Println("for 循环外部")

	<-stopper // 从channel 中获取数据

	//time.Sleep(time.Second * 30)

}
