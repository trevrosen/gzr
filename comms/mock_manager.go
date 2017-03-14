package comms

type MockManager struct {
	OnBuild func(...string) error
	OnPush  func(string) error
}

func (mock *MockManager) Build(args ...string) error {
	return mock.OnBuild(args...)
}

func (mock *MockManager) Push(name string) error {
	return mock.OnPush(name)
}
