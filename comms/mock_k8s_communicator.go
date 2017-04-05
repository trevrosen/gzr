package comms

type MockK8sCommunicator struct {
	OnGetDeployment    func(string) (*GzrDeployment, error)
	OnListDeployments  func() (*GzrDeploymentList, error)
	OnUpdateDeployment func(*DeploymentContainerInfo) (*GzrDeployment, error)

	namespace string
}

func (mock *MockK8sCommunicator) GetDeployment(deploymentName string) (*GzrDeployment, error) {
	return mock.OnGetDeployment(deploymentName)
}

func (mock *MockK8sCommunicator) ListDeployments() (*GzrDeploymentList, error) {
	return mock.OnListDeployments()
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
