package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/varshaprasad96/custom-crd-operator/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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

	crdConfig := *cfg
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	res := v1alpha1.ProjectList{}
	err = exampleRestClient.Get().Resource("projects").Do(context.TODO()).Into(&res)

	if err != nil {
		panic(err)
	}
	// clientSet, err := clientV1alpha1.NewForConfig(cfg)
	// if err != nil {
	// 	panic(err)
	// }

	// res, err := clientSet.Projects("default").List(metav1.ListOptions{}, context.TODO())
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("project found *****************")
	fmt.Println(res.GetResourceVersion())
	fmt.Println(res.GroupVersionKind().String())
	fmt.Printf("projects found: %+v\n", res)

}
