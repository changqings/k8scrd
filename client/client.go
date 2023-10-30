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

// get dynamic client with Context
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

// get dyn client for use dynamic Interface
func GetDynamicClient() dynamic.Interface {
	config := GetRestConfig()
	// creates the clientset, default behavor

	dynaClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return dynaClient
}

// k8s api client with set-context = contextName
// configPath equal merged kubeconfigs, example dev, prod, test
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

// k8s api client
func GetClient() *kubernetes.Clientset {
	config := GetRestConfig()
	// creates the clientset, default behavor

	clientset, err := kubernetes.NewForConfig(config)
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
