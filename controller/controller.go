package controller

import (
	"github.com/openshift/library-go/pkg/operator/events"
	"k8s.io/client-go/kubernetes"
)

type TestController struct {
	name       string
	kubeclinet kubernetes.Interface
	recorder   events.Recorder
}
