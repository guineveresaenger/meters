package main

import (
	"github.com/guineveresaenger/meters/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	provLibrary "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"

)

type iceCreamAdapter struct {
	cmd.AdapterBase
}

func main() {
	adapter := &iceCreamAdapter{}

	iceCreamProvider := adapter.makeProviderOrDie()

	adapter.WithCustomMetrics(iceCreamProvider)
	klog.Infof("Starting adapter!")
	if err := adapter.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}

func (ica *iceCreamAdapter) makeProviderOrDie() provLibrary.CustomMetricsProvider {
	client, err := ica.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := ica.RESTMapper()
	if err != nil {
		klog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return provider.NewProvider(mapper, client)
}