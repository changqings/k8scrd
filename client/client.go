package client

import (
	"log/slog"
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

type Client struct {
	// kubeconfig path
	KubeConfig string
	// REST config, use for other resources
	RestConfig *rest.Config
	// kube clientset
	KubeClient *kubernetes.Clientset
}

func NewClient() (*Client, error) {

	kubeConfig := GetKubeConfig()

	restConfig, err := GetRestConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	kubeClient, err := GetClient(restConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		KubeConfig: kubeConfig,
		RestConfig: restConfig,
		KubeClient: kubeClient,
	}, nil
}

// get istio client
func (c *Client) GetIstioClient() (*istioVersioned.Clientset, error) {
	return istioVersioned.NewForConfig(c.RestConfig)
}

// get runtime client
func (c *Client) GetRuntimeClient(r *runtime.Scheme) (runtimeclient.Client, error) {
	// creates the runtime client with  scheme
	client, err := runtimeclient.New(c.RestConfig, runtimeclient.Options{
		Scheme: r,
	})
	if err != nil {
		slog.Error("NewRuntimeClient", "msg", err)
		return nil, err
	}
	return client, nil
}

// get dynamic client with Context
func (c *Client) GetDynamicClientWithContext(contextName string) (*dynamic.DynamicClient, error) {

	var err error
	var restConfig *rest.Config

	if fileIsExist(c.KubeConfig) && len(contextName) > 0 {
		configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: c.KubeConfig}
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
		if err != nil {
			slog.Error("Switch kubeconfig context err", "msg", err)
			return nil, err
		}
	} else {
		restConfig = c.RestConfig
	}

	return dynamic.NewForConfig(restConfig)
}

// get dyn client for use dynamic DynamicClient
func (c *Client) GetDynamicClient() (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(c.RestConfig)
}

// k8s api client with set-context = contextName
// configPath equal merged kubeconfigs, example dev, prod, test
func (c *Client) GetClientSetWithContext(contextName string) (*kubernetes.Clientset, error) {
	var restConfig *rest.Config
	var err error

	if fileIsExist(c.KubeConfig) && len(contextName) > 0 {
		configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: c.KubeConfig}
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
		if err != nil {
			klog.Errorf("Switch kubeconfig context err:  %v", err)
			slog.Error("Switch kubeconfig context err", "msg", err)
			return nil, err
		}
	} else {
		restConfig = c.RestConfig
	}

	return kubernetes.NewForConfig(restConfig)
}

// k8s clientset
func GetClient(restConfig *rest.Config) (*kubernetes.Clientset, error) {

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		slog.Error("NewForConfig", "msg", err)
		return nil, err
	}
	return clientSet, nil
}

// first use kubeconfig from configPath
// second use env.KUBECONFIG
// third use $HOME/.kube/config
// default return empty string
func GetKubeConfig(configPath ...string) string {

	homeDir := homedir.HomeDir()
	kubeconfigEnv := os.Getenv("KUBECONFIG")

	var kubeconfig string
	switch {
	case len(configPath) > 0:
		kubeconfig = configPath[0]
	case len(kubeconfigEnv) > 0:
		kubeconfig = filepath.Join(kubeconfigEnv)
	case len(homeDir) > 0:
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
	default:
		kubeconfig = ""
	}

	return kubeconfig
}

// get rest config, if you want use for other resources
// Example: istioClient
func GetRestConfig(kubeconfig string) (*rest.Config, error) {

	var resConfig *rest.Config
	var err error

	if fileIsExist(kubeconfig) {
		resConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			slog.Error("buildConfigFromFlags", "msg", err)
			return nil, err
		}
	} else {
		// creates the in-cluster config
		resConfig, err = rest.InClusterConfig()
		if err != nil {
			slog.Error("InClusterConfig", "msg", err)
			return nil, err
		}
	}

	return resConfig, nil
}

func fileIsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
