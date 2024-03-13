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

// get istio client
func GetIstioClient() *istioVersioned.Clientset {
	return istioVersioned.NewForConfigOrDie(GetRestConfig())
}

// get runtime client
func GetRuntimeClient(r *runtime.Scheme) runtimeclient.Client {
	resConfig := GetRestConfig()

	// creates the clientset, default behavor
	client, err := runtimeclient.New(resConfig, runtimeclient.Options{
		Scheme: r,
	})
	if err != nil {
		panic(err.Error())
	}
	return client
}

// get dynamic client with Context
func GetDynamicClientWithContext(contextName string) *dynamic.DynamicClient {
	var resConfig *rest.Config

	var err error

	kubeconfig := GetKubeConfig()
	if fileIsExist(kubeconfig) {
		resConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if contextName != "" {
			configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
			configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
			resConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()

			if err != nil {
				klog.Errorf("Switch kubeconfig context err: ", err)
			}
		}
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset, default behavor
	dynaClient, err := dynamic.NewForConfig(resConfig)
	if err != nil {
		panic(err.Error())
	}
	return dynaClient
}

// get dyn client for use dynamic DynamicClient
func GetDynamicClient() *dynamic.DynamicClient {
	resConfig := GetRestConfig()
	// creates the clientset, default behavor

	dynClient, err := dynamic.NewForConfig(resConfig)
	if err != nil {
		panic(err.Error())
	}
	return dynClient
}

// k8s api client with set-context = contextName
// configPath equal merged kubeconfigs, example dev, prod, test
func GetClientWithContext(contextName string, configPath string) *kubernetes.Clientset {
	var resConfig *rest.Config
	var err error

	kubeconfig := GetKubeConfig(configPath)
	if fileIsExist(kubeconfig) {
		resConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if contextName != "" {
			configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
			configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
			resConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()

			if err != nil {
				klog.Errorf("Switch kubeconfig context err: ", err)
			}
		}
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset, default behavor

	clientset, err := kubernetes.NewForConfig(resConfig)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

// k8s api client
func GetClient() *kubernetes.Clientset {
	resConfig := GetRestConfig()
	// creates the clientset, default behavor

	clientset, err := kubernetes.NewForConfig(resConfig)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

// if configPath existed, use it first
// or use ENV KUBECONFIG
func GetKubeConfig(configPath ...string) string {

	var kubeconfig string
	homeDir := homedir.HomeDir()
	ENV_KUBUCONFIG := os.Getenv("KUBECONFIG")

	switch {
	case len(configPath) == 1:
		kubeconfig = configPath[0]
	case len(ENV_KUBUCONFIG) > 0:
		kubeconfig = filepath.Join(ENV_KUBUCONFIG)
	case len(homeDir) > 0:
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
	default:
		kubeconfig = ""
	}

	return kubeconfig
}

// get rest config, if you want use for other resources
// Example: istioClient
func GetRestConfig() *rest.Config {

	kubeconfig := GetKubeConfig()
	var resConfig *rest.Config
	var err error

	if fileIsExist(kubeconfig) {
		resConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		// creates the in-cluster config
		resConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	return resConfig
}

func fileIsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
