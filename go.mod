module github.com/guineveresaenger/meters

go 1.13

require (
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190817021527-637fc595d17a
	k8s.io/klog v0.3.1
	k8s.io/metrics v0.0.0-20190817023635-63ee757b2e8b
)

replace github.com/kubernetes-incubator/custom-metrics-apiserver => ../../../kubernetes-incubator/custom-metrics-apiserver
