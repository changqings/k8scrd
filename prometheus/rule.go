package prometheus

import (
	"context"
	"fmt"

	"github.com/Tsingshen/k8s-crd/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetP8sRule(name, ns string) error {
	dynaClient := client.GetDynamicClient()
	// get gvr prometheusrulers
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.coreos.com",
		Version:  "v1",
		Resource: "prometheusrules",
	}

	unStructOjb, err := dynaClient.Resource(gvr).Namespace(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range unStructOjb.Items {
		fmt.Printf("%s\n", d.GetName())
	}

	return nil

}
