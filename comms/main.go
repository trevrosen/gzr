package comms

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func EstablishK8sConnection(kubeconfig string) error {
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentList, err := clientset.ExtensionsV1beta1().Deployments("default").List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, deployment := range deploymentList.Items {
		var displayLabel string

		// Deployments either have an app or run label so far. I'm not sure how they get one vs the other.
		// -- The hello-minikube deploy is the "run" one
		// -- The nginx deploy is the "app" one
		labels := deployment.Spec.Template.Labels
		run := labels["run"]
		app := labels["app"]

		if app == "" {
			displayLabel = run
		} else {
			displayLabel = app
		}

		fmt.Printf("Deployment: %s\n", displayLabel)
		fmt.Printf("-- selector: %s\n", deployment.Spec.Selector)

	}

	return err
}
