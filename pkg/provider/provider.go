package provider

import (
    "fmt"
    "time"

    "k8s.io/klog"
    apierr "k8s.io/apimachinery/pkg/api/errors"
    apimeta "k8s.io/apimachinery/pkg/api/meta"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/dynamic"
    "k8s.io/metrics/pkg/apis/custom_metrics"

    "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
    "k8s.io/metrics/pkg/apis/external_metrics"
)

type CustomMetricsProvider interface {
    ListAllMetrics() []CustomMetricInfo

    GetMetricByName(name types.NamespacedName, info CustomMetricInfo) (*custom_metrics.MetricValue, error)
    GetMetricBySelector(namespace string, selector labels.Selector, info CustomMetricInfo) (*custom_metrics.MetricValueList, error)
}

func (p *yourProvider) ListAllMetrics() []provider.CustomMetricInfo {
    return []provider.CustomMetricInfo{
        // these are mostly arbitrary examples
        {
            GroupResource: schema.GroupResource{Group: "", Resource: "pods"},
            Metric:        "packets-per-second",
            Namespaced:    true,
        },
        {
            GroupResource: schema.GroupResource{Group: "", Resource: "services"},
            Metric:        "connections-per-second",
            Namespaced:    true,
        },
        {
            GroupResource: schema.GroupResource{Group: "", Resource: "namespaces"},
            Metric:        "work-queue-length",
            Namespaced:    false,
        },
    }
}

type yourProvider struct {
    client dynamic.Interface
    mapper apimeta.RESTMapper

    // just increment values when they're requested
    values map[provider.CustomMetricInfo]int64
}

// valueFor fetches a value from the fake list and increments it.
func (p *yourProvider) valueFor(info provider.CustomMetricInfo) (int64, error) {
    // normalize the value so that you treat plural resources and singular
    // resources the same (e.g. pods vs pod)
    info, _, err := info.Normalized(p.mapper)
    if err != nil {
        return 0, err
    }

    value := p.values[info]
    value += 1
    p.values[info] = value

    return value, nil
}

// metricFor constructs a result for a single metric value.
func (p *testingProvider) metricFor(value int64, name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
    // construct a reference referring to the described object
    objRef, err := helpers.ReferenceFor(p.mapper, name, info)
    if err != nil {
        return nil, err
    }

    return &custom_metrics.MetricValue{
        DescribedObject: objRef,
        MetricName:      info.Metric,
        // you'll want to use the actual timestamp in a real adapter
        Timestamp:       metav1.Time{time.Now()},
        Value:           *resource.NewMilliQuantity(value*100, resource.DecimalSI),
    }, nil
}

func (p *yourProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
    value, err := p.valueFor(info)
    if err != nil {
        return nil, err
    }
    return p.metricFor(value, name, info)
}

func (p *yourProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo) (*custom_metrics.MetricValueList, error) {
    totalValue, err := p.valueFor(info)
    if err != nil {
        return nil, err
    }

    names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
    if err != nil {
        return nil, err
    }

    res := make([]custom_metrics.MetricValue, len(names))
    for i, name := range names {
        // in a real adapter, you might want to consider pre-computing the
        // object reference created in metricFor, instead of recomputing it
        // for each object.
        value, err := p.metricFor(100*totalValue/int64(len(res)), types.NamespacedName{Namespace: namespace, Name: name}, info)
        if err != nil {
            return nil, err
        }
        res[i] = *value
    }

    return &custom_metrics.MetricValueList{
        Items: res,
    }, nil
}

