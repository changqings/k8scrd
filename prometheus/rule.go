package prometheus

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	p8smonitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func GetP8sRule(dI dynamic.Interface, name, ns string) error {
	// get gvr prometheusrulers
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.coreos.com",
		Version:  "v1",
		Resource: "prometheusrules",
	}

	unStructObj, err := dI.Resource(gvr).Namespace(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range unStructObj.Items {
		fmt.Printf("%s\n", d.GetName())
	}

	return nil

}

func GetP8sRuleList(rClient client.Client) *p8smonitorv1.PrometheusRuleList {

	// ruleScheme := runtime.NewScheme()
	// p8smonitorv1.AddToScheme(ruleScheme)

	// rClient := k8scrdclient.GetRuntimeClient(ruleScheme)

	p8sRule := &p8smonitorv1.PrometheusRuleList{}
	err := rClient.List(context.TODO(), p8sRule, client.InNamespace("prometheus"))

	if err != nil {
		panic(err)
	}

	return p8sRule

}
