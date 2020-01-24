package provider


import (
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type iceCreamMetricsProvider struct {

}

func (ic *iceCreamMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return []provider.CustomMetricInfo{
		// these are mostly delicious examples
		{
			GroupResource: schema.GroupResource{Group: "", Resource: "pods"},
			Metric:        "scoops-per-second",
			Namespaced:    true,
		},
		{
			GroupResource: schema.GroupResource{Group: "", Resource: "services"},
			Metric:        "sprinkles-per-second",
			Namespaced:    true,
		},
		{
			GroupResource: schema.GroupResource{Group: "", Resource: "namespaces"},
			Metric:        "whipped-cream",
			Namespaced:    false,
		},
	}
}