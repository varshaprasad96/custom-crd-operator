package main

import (
	"context"
	"flag"
	"time"

	"github.com/openshift/library-go/pkg/operator/events"
	log "github.com/sirupsen/logrus"
	"github.com/varshaprasad96/custom-crd-operator/controller"
	"github.com/varshaprasad96/custom-crd-operator/pkg/apis/example/v1alpha1"
	operatorversionedclient "github.com/varshaprasad96/custom-crd-operator/pkg/generated/clientset/versioned"
	clientV1alpha1 "github.com/varshaprasad96/custom-crd-operator/pkg/generated/clientset/versioned/typed/example/v1alpha1"
	opInformer "github.com/varshaprasad96/custom-crd-operator/pkg/generated/informers/externalversions"
	coreinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

// Pass kubeconfig as a flag or specify a flag.
func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	ctx := context.TODO()
	var cfg *rest.Config
	var err error

	// Hardcoding the kubeconfig filepath for now. Accept it from flag.
	cfg, err = clientcmd.BuildConfigFromFlags("", "config")
	if err != nil {
		panic(err)
	}

	// Register custom resource to the scheme.
	v1alpha1.AddToScheme(scheme.Scheme)

	// Create a config to work with the custom resource.
	clientSet, err := clientV1alpha1.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// create a versioned client which will be used to create informers in turn.
	operatorConfigClient, err := operatorversionedclient.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// create kubeclient to handle other resources like deployment.
	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// create an informer for memcached resource.
	memcachedInformer := opInformer.NewSharedInformerFactoryWithOptions(
		operatorConfigClient,
		time.Minute,
	)

	// use coreInformer to set up an informer for the deployment object.
	coreInformerFactory := coreinformers.NewSharedInformerFactory(kubeclient, 0)
	testController := controller.NewTestController("memcached-sample", clientSet, kubeclient, coreInformerFactory.Apps().V1().Deployments(), events.NewInMemoryRecorder("memcached"), memcachedInformer.Example().V1alpha1().Memcacheds(), "default")

	log.Infof("starting informers")
	// Start the informers to make sure their caches are in sync and are updated periodically.
	for _, informer := range []interface {
		Start(stopCh <-chan struct{})
	}{
		coreInformerFactory,
		memcachedInformer,
	} {
		informer.Start(ctx.Done())
	}

	log.Infof("starting controllers")
	// Start and run the controller.
	for _, controllerint := range []interface {
		Run(ctx context.Context, workers int)
	}{
		testController,
	} {
		go controllerint.Run(ctx, 1)
	}

	<-ctx.Done()
	return
}
