package main

import (
	"context"
	"fmt"

	crdClient "github.com/Tsingshen/k8scrd/Client"
	"github.com/Tsingshen/k8scrd/crd"
	"github.com/Tsingshen/k8scrd/prometheus"
	p8smonitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	dI := crdClient.GetDynamicClient()
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
	if err := prometheus.GetP8sRule(dI, "", "prometheus"); err != nil {
		fmt.Printf("%v\n", err)
	}

	// get vs crds
	fmt.Printf("get vs with dynamic client")
	vsGvr := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1beta1",
		Resource: "virtualservices",
	}

	vs := &crd.Crds{
		Namespace: "shencq",
		Gvr:       vsGvr,
		Client:    dI,
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

}