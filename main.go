package main

import (
	"flag"
	"fmt"

	"github.com/varshaprasad96/custom-crd-operator/api/types/v1alpha1"
	clientV1alpha1 "github.com/varshaprasad96/custom-crd-operator/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	res, err := clientSet.Projects("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("project found *****************")
	fmt.Println(res.GetResourceVersion())
	fmt.Println(res.GroupVersionKind().String())
	fmt.Printf("projects found: %+v\n", res)

}
