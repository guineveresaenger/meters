package provider


import (
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider/helpers"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"time"
)

type iceCreamMetricsProvider struct {
	mapper apimeta.RESTMapper

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

func (ic *iceCreamMetricsProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	var val int64
	val = 1

	objRef, err := helpers.ReferenceFor(ic.mapper, name, info)
	if err != nil {
		return nil, err
	}

	cmv := &custom_metrics.MetricValue{
		DescribedObject: objRef,
		Metric: custom_metrics.MetricIdentifier{
			Name: info.Metric,
		},
		Timestamp: metav1.Time{time.Now()},
		Value: *resource.NewQuantity(val, resource.DecimalSI),

	}
	return cmv, nil
}


