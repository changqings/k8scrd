package main

import (
	"context"
	"fmt"
	"log"

	"github.com/changqings/k8scrd/client"
	"github.com/changqings/k8scrd/crd"
	"github.com/changqings/k8scrd/prometheus"
	p8smonitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	runtimeshceme = runtime.NewScheme()
)

func main() {
	fmt.Println("Runnig main ...")

	// take prometheusRules as example
	// use k8s dyn client
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	dynClient, err := client.GetDynamicClient()
	if err != nil {
		log.Fatal(err)
	}

	p8smonitorv1.AddToScheme(runtimeshceme)
	rClient, err := client.GetRuntimeClient(runtimeshceme)
	if err != nil {
		log.Fatal(err)
	}

	// list prometheusRules
	list := prometheus.GetP8sRuleList(rClient)
	for _, v := range list.Items {
		fmt.Printf("%s\n", v.Name)
	}
	fmt.Println("get p8srules with dynamic client")

	if err := prometheus.GetP8sRule(dynClient, "", "prometheus"); err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("get vs with dynamic client")

	// istio virtualService use dynamic client
	vsGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1beta1",
		Resource: "virtualservices",
	}

	vsName := "nginx-vs"
	vs := &crd.CrdClient{
		Namespace: "shencq",
		GVR:       vsGVR,
		Client:    dynClient,
	}

	// vs get
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

	// vs create wiht unstructured obj, best to use istio api client, but here just for demo
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

	_, err = vs.Create(context.Background(), vsUnObj, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create vs err: %v\n", err)
	}

}
