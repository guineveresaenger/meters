package provider


import (
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider/helpers"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"time"
)

type iceCreamMetricsProvider struct {
	mapper apimeta.RESTMapper
	client dynamic.Interface

}

func NewProvider(m apimeta.RESTMapper, c dynamic.Interface) provider.CustomMetricsProvider {
	return &iceCreamMetricsProvider{
		mapper: m,
		client: c,
	}
}

func (ic *iceCreamMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return []provider.CustomMetricInfo{
		// these are mostly delicious examples
		// this needs to be replaced by reading in all the actual metrics
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

func (ic *iceCreamMetricsProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo, metricsSelector labels.Selector) (*custom_metrics.MetricValue, error) {
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

func (ic *iceCreamMetricsProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricsSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	names, err := helpers.ListObjectNames(ic.mapper, ic.client, namespace, selector, info)
	if err != nil {
		return nil, err
	}
	list := make([]custom_metrics.MetricValue, len(names))

	for i, name := range names {
		val, _ := ic.GetMetricByName(types.NamespacedName{Namespace: namespace, Name: name}, info, metricsSelector)
		list[i] = *val
	}
	return &custom_metrics.MetricValueList{
		Items:    list,
	}, nil
}

