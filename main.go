package main

import (
	"context"
	"fmt"

	crdClient "github.com/changqings/k8scrd/client"
	"github.com/changqings/k8scrd/crd"
	"github.com/changqings/k8scrd/prometheus"
	p8smonitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	dynClient := crdClient.GetDynamicClient()
	fmt.Println("Runnig main ...")

	// with runtime client
	fmt.Println("get p8srules with rest client")
	ruleScheme := runtime.NewScheme()
	p8smonitorv1.AddToScheme(ruleScheme)
	rClient := crdClient.GetRuntimeClient(ruleScheme)
	list := prometheus.GetP8sRuleList(rClient)
	for _, v := range list.Items {
		fmt.Printf("%s\n", v.Name)
	}

	fmt.Println("get p8srules with dynamic client")
	// with dynamic client
	if err := prometheus.GetP8sRule(dynClient, "", "prometheus"); err != nil {
		fmt.Printf("%v\n", err)
	}

	// get vs crds
	fmt.Printf("get vs with dynamic client")
	vsGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1beta1",
		Resource: "virtualservices",
	}

	vs := &crd.Crds{
		Namespace: "shencq",
		GVR:       vsGVR,
		Client:    dynClient,
	}

	vsName := "nginx-vs"

	unobj, err := vs.Get(context.Background(), vsName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("get vs unobj error %v\n", err)
		return
	}

	fmt.Printf("get vs %s/%s ok\n", vs.Namespace, unobj.GetName())
	vsSpec, ok, _ := unstructured.NestedMap(unobj.Object, "spec")
	if ok {
		fmt.Println("get vsSpec ok, and k/v =")
		for k, v := range vsSpec {
			fmt.Printf("%s: %v\n", k, v)

		}
	}

	// create vs

	vsUnObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.istio.io/v1beta1",
			"kind":       "VirtualService",
			"metadata": map[string]string{
				"name":      "test-crd-vs",
				"namespace": "shencq",
			},
			"spec": map[string]interface{}{
				"gateways": []string{"mesh"},
				"hosts":    []string{"test.abc.com"},
				"http": []map[string]interface{}{
					{
						"name": "test-crd" + "-stable",
						"route": []map[string]interface{}{
							{
								"destination": map[string]string{
									"host":   "test-crd-vs.shencq.svc.cluster.local",
									"subset": "stable",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err1 := vs.Create(context.Background(), vsUnObj, metav1.CreateOptions{})
	if err1 != nil {
		panic(err1)
	}

}
