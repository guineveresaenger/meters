module github.com/guineveresaenger/meters

go 1.13

require (
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0
	k8s.io/apimachinery v0.17.2 // indirect
)

replace github.com/kubernetes-incubator/custom-metrics-apiserver => ../../../kubernetes-incubator/custom-metrics-apiserver
