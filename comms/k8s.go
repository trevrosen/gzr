package comms

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

type GzrDeployment v1beta1.Deployment

type Serializer interface {
	// SerializeForCLI writes templatized information to the provided io.Writer
	SerializeForCLI(io.Writer) error
	// SerializeForWeb kicks out JSON as a byte slice
	SerializeForWeb() ([]byte, error)
}

// K8sCommunicator defines an interface for retrieving data from a k8s cluster
type K8sCommunicator interface {
	// Deployments returns the list of Deployments in the cluster
	Deployments() ([]v1beta1.Deployment, error)
}

// K8sConnection implements the K8sCommunicator interface and holds a live connection to a k8s cluster
type K8sConnection struct {
	// clientset is a collection of Kubernetes API clients
	clientset *kubernetes.Clientset
}

// NewK8sConnection returns a K8sConnection with an active v1.Clientset.
//   - assumes that $HOME/.kube/config contains a legit Kubernetes config for an healthy k8s cluster.
//   - panics if the configuration can't be used to connect to a k8s cluster.
func NewK8sConnection() (*K8sConnection, error) {
	var k *K8sConnection
	kubeconfig := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return k, err
	}

	k = &K8sConnection{
		clientset: clientset,
	}

	return k, nil
}

// DeploymentByName returns a k8s v1.Deployment
//func (k *K8sConnection) DeploymentByName(name string) *v1.Deployment {

//}

// Deployments returns the active k8s Deployments
func (k *K8sConnection) Deployments(namespace string) ([]GzrDeployment, error) {
	var deployments []GzrDeployment
	deploymentList, err := k.clientset.ExtensionsV1beta1().Deployments(namespace).List(v1.ListOptions{})
	if err != nil {
		return deployments, err
	}

	for _, deployment := range deploymentList.Items {
		deployments = append(deployments, GzrDeployment(deployment))
	}

	return deployments, nil
}

// SerializeForCLI takes an io.Writer and writes templatized data to it representing a Deployment
func (d GzrDeployment) SerializeForCLI(wr io.Writer) error {
	return d.cliTemplate().Execute(wr, d)
}

// cliTemplate returns the template that will be used for serializing Deployment data to the CLI
func (d GzrDeployment) cliTemplate() *template.Template {
	t := template.New("Deployment CLI")
	t, _ = t.Parse(`-------------------------
Deployment: {{.ObjectMeta.Name}}
  - containers: {{range .Spec.Template.Spec.Containers}}
    --name:  {{.Name}}
    --image: {{.Image}}
{{end}}
`)
	return t
}
