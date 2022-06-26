package prometheus

import (
	"context"
	"fmt"

	k8scrdclient "github.com/Tsingshen/k8scrd/client"
	"sigs.k8s.io/controller-runtime/pkg/client"

	p8smonitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetP8sRule(name, ns string) error {
	dynaClient := k8scrdclient.GetDynamicClient()
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

func GetP8sRuleList() *p8smonitorv1.PrometheusRuleList {

	ruleScheme := runtime.NewScheme()
	p8smonitorv1.AddToScheme(ruleScheme)

	rClient := k8scrdclient.GetRuntimeClient(ruleScheme)

	p8sRule := &p8smonitorv1.PrometheusRuleList{}
	err := rClient.List(context.TODO(), p8sRule, client.InNamespace("prometheus"))

	if err != nil {
		panic(err)
	}

	return p8sRule

}
