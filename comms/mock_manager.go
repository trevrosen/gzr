package comms

type MockManager struct {
	OnBuild func(...string) error
	OnPush  func(string) error
}

// NewDefaultMockManager returns an initialized MockManager with no-ops for all methods
func NewDefaultMockManager() ImageManager {
	return &MockManager{
		OnBuild: func(_ ...string) error { return nil },
		OnPush:  func(_ string) error { return nil },
	}
}

func (mock *MockManager) Build(args ...string) error {
	return mock.OnBuild(args...)
}

func (mock *MockManager) Push(name string) error {
	return mock.OnPush(name)
}
