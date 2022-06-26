package client

import (
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetRuntimeClient(r *runtime.Scheme) runtimeclient.Client {
	var kubeconfig string
	var config *rest.Config

	var err error
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = filepath.Join(os.Getenv("KUBECONFIG"))
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset, default behavor

	client, err := runtimeclient.New(config, runtimeclient.Options{
		Scheme: r,
	})
	if err != nil {
		panic(err.Error())
	}
	return client
}

func GetDynamicClient() dynamic.Interface {
	var kubeconfig string
	var config *rest.Config

	var err error
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = filepath.Join(os.Getenv("KUBECONFIG"))
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
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

func GetClient() *kubernetes.Clientset {
	var kubeconfig string
	var config *rest.Config

	var err error
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = filepath.Join(os.Getenv("KUBECONFIG"))
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	if fileExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
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

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
