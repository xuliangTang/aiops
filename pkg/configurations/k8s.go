package configurations

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

type K8sConfig struct{}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (*K8sConfig) K8sRestConfigDefault() *rest.Config {
	// 取用户目录   Linux： ~   /home/xxx
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// 从 ~/.kube/config取
	defaultConfigPath := fmt.Sprintf("%s/.kube/config", home)

	config, err := clientcmd.BuildConfigFromFlags("", defaultConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true
	return config
}

// InitDynamicClient 初始化动态客户端
func (ths *K8sConfig) InitDynamicClient() dynamic.Interface {
	client, err := dynamic.NewForConfig(ths.K8sRestConfigDefault())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// RestMapper 获取所有apiGroupResource
// 这个要缓存起来，不然反复从k8s api获取会比较慢
func (ths *K8sConfig) RestMapper() *meta.RESTMapper {
	gr, err := restmapper.GetAPIGroupResources(ths.InitClient().Discovery())
	if err != nil {
		log.Fatal(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	return &mapper
}

// InitClient 初始化clientset
func (ths *K8sConfig) InitClient() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(ths.K8sRestConfigDefault())
	if err != nil {
		log.Fatal(err)
	}

	return c
}

// InitWatchFactory 初始化一个动态客户端
func (ths *K8sConfig) InitWatchFactory() dynamicinformer.DynamicSharedInformerFactory {
	dynclient := ths.InitDynamicClient() // 取出动态客户端

	fact := dynamicinformer.NewDynamicSharedInformerFactory(dynclient, 0)
	// 临时的
	fact.ForResource(schema.
		GroupVersionResource{Version: "v1", Resource: "namespaces"}).
		Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	fact.ForResource(schema.
		GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).
		Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	fact.Start(wait.NeverStop)
	return fact
}
