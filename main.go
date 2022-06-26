package main

import (
	"fmt"

	"github.com/Tsingshen/k8scrd/prometheus"
)

func main() {
	fmt.Println("Runnig main ...")
	// with runtime client
	list := prometheus.GetP8sRuleList()
	for _, v := range list.Items {
		fmt.Printf("%s\n", v.Name)
	}

	fmt.Println("------")
	// with dynamic client
	if err := prometheus.GetP8sRule("", "prometheus"); err != nil {
		fmt.Printf("%v\n", err)
	}

}
