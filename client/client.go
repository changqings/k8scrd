package client

import (
	"os"
	"path/filepath"

	istioVersioned "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetIstioClient() *istioVersioned.Clientset {
	return istioVersioned.NewForConfigOrDie(GetRestConfig())
}

func GetRuntimeClient(r *runtime.Scheme) runtimeclient.Client {
	config := GetRestConfig()
	// creates the clientset, default behavor

	client, err := runtimeclient.New(config, runtimeclient.Options{
		Scheme: r,
	})
	if err != nil {
		panic(err.Error())
	}
	return client
}

func GetDynamicClientWithContext(contextName string) dynamic.Interface {
	var config *rest.Config

	var err error

	kubeconfig := GetKubeConfig()
	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if contextName != "" {
			configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
			configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()

			if err != nil {
				klog.Errorf("Switch kubeconfig context err: ", err)
			}
		}
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset, default behavor

	dynaClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return dynaClient
}
func GetDynamicClient() dynamic.Interface {
	config := GetRestConfig()
	// creates the clientset, default behavor

	dynaClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return dynaClient
}

func GetClientWithContext(contextName string, configPath string) *kubernetes.Clientset {
	var config *rest.Config
	var err error

	kubeconfig := GetKubeConfig(configPath)
	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if contextName != "" {
			configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
			configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()

			if err != nil {
				klog.Errorf("Switch kubeconfig context err: ", err)
			}
		}
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset, default behavor

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
func GetClient() *kubernetes.Clientset {
	config := GetRestConfig()
	// creates the clientset, default behavor

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func GetKubeConfig(configPath ...string) string {
	var kubeconfig string

	if len(configPath) == 1 {
		kubeconfig = configPath[0]
	} else if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = filepath.Join(os.Getenv("KUBECONFIG"))
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	return kubeconfig

}

func GetRestConfig() *rest.Config {

	kubeconfig := GetKubeConfig()
	var config *rest.Config
	var err error

	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	return config
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
