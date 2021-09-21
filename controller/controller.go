package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/varshaprasad96/custom-crd-operator/api/types/v1alpha1"
	clientV1alpha1 "github.com/varshaprasad96/custom-crd-operator/clientset/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
)

type TestController struct {
	name           string
	operatorClient *clientV1alpha1.ExampleV1Alpha1Client
	kubeclient     kubernetes.Interface
	deployInformer appsinformersv1.DeploymentInformer
	recorder       events.Recorder
	namespace      string
}

func NewTestController(name string,
	operatorClient *clientV1alpha1.ExampleV1Alpha1Client,
	kubeclient kubernetes.Interface,
	deployInformer appsinformersv1.DeploymentInformer,
	recorder events.Recorder,
	ns string) factory.Controller {
	c := &TestController{
		name:           name,
		operatorClient: operatorClient,
		kubeclient:     kubeclient,
		deployInformer: deployInformer,
		recorder:       recorder,
		namespace:      ns,
	}

	return factory.New().WithInformers(deployInformer.Informer()).WithSync(c.sync).ResyncEvery(time.Minute).ToController(c.name, recorder.WithComponentSuffix(strings.ToLower(name)+"-deployment-controller-"))
}

func (c *TestController) sync(ctx context.Context, syncContext factory.SyncContext) error {

	fmt.Println("reconciling************")
	project, err := c.operatorClient.Projects(c.namespace).Get(c.name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("project resource not found. Ignoring and reconciling again since object maybe deleted")
			return nil
		}
		fmt.Println("failed to get project")
		return nil
	}

	found, err := c.kubeclient.AppsV1().Deployments(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		dep := c.deploymentForProject(project)
		fmt.Println("creating new deployment")
		_, err := c.kubeclient.AppsV1().Deployments(c.namespace).Create(ctx, dep, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return nil
	} else if err != nil {
		fmt.Println(err)
	}

	size := project.Spec.Replicas
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		_, err := c.kubeclient.AppsV1().Deployments(c.namespace).Update(ctx, found, v1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	return nil
}

func (r *TestController) deploymentForProject(m *v1alpha1.Project) *appsv1.Deployment {
	ls := labelsForMemcached(m.Name)
	replica := m.Spec.Replicas

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "memcached:1.4.36-alpine",
						Name:    "projects",
						Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "projects",
						}},
					}},
				},
			},
		},
	}
	return dep
}

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForMemcached(name string) map[string]string {
	return map[string]string{"app": "project", "project_cr": name}
}
