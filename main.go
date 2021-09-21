package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/varshaprasad96/custom-crd-operator/controller"
	"github.com/varshaprasad96/custom-crd-operator/pkg/apis/example.com/v1alpha1"
	operatorversionedclient "github.com/varshaprasad96/custom-crd-operator/pkg/generated/clientset/versioned"
	clientV1alpha1 "github.com/varshaprasad96/custom-crd-operator/pkg/generated/clientset/versioned/typed/example.com/v1alpha1"
	opInformer "github.com/varshaprasad96/custom-crd-operator/pkg/generated/informers/externalversions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var cfg *rest.Config
	var err error

	fmt.Println("using configuration from '%s'", kubeconfig)
	cfg, err = clientcmd.BuildConfigFromFlags("", "config")

	if err != nil {
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	operatorConfigClient, err := operatorversionedclient.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	res, err := clientSet.Projects("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("project found *****************")
	fmt.Println(res.GetResourceVersion())
	fmt.Println(res.GroupVersionKind().String())
	fmt.Printf("projects found: %+v\n", res)

	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	projectInformer := opInformer.NewSharedInformerFactoryWithOptions(
		operatorConfigClient,
		time.Minute,
	)

	ctx := context.TODO()
	coreInformerFactory := coreinformers.NewSharedInformerFactory(kubeclient, 0)
	testController := controller.NewTestController("example-project", clientSet, kubeclient, coreInformerFactory.Apps().V1().Deployments(), events.NewInMemoryRecorder("memcached"), projectInformer.Example().V1alpha1().Projects(), "default")

	for _, informer := range []interface {
		Start(stopCh <-chan struct{})
	}{
		coreInformerFactory,
		projectInformer,
	} {
		informer.Start(ctx.Done())
	}

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
