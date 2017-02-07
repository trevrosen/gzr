package comms

type MockK8sCommunicator struct {
	OnGetDeployment    func(string, string) (*GzrDeployment, error)
	OnListDeployments  func(string) (*GzrDeploymentList, error)
	OnUpdateDeployment func(*DeploymentContainerInfo) (*GzrDeployment, error)

	namespace string
}

func (mock *MockK8sCommunicator) GetDeployment(namespace string, deploymentName string) (*GzrDeployment, error) {
	return mock.OnGetDeployment(namespace, deploymentName)
}

func (mock *MockK8sCommunicator) ListDeployments(namespace string) (*GzrDeploymentList, error) {
	return mock.OnListDeployments(namespace)
}

func (mock *MockK8sCommunicator) UpdateDeployment(dci *DeploymentContainerInfo) (*GzrDeployment, error) {
	return mock.OnUpdateDeployment(dci)
}

func (mock *MockK8sCommunicator) GetNamespace() string {
	if mock.namespace == "" {
		return "default"
	}
	return mock.namespace
}
