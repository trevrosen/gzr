package comms

type MockBuilder struct {
	OnBuild func(...string) error
	OnPush  func(string) error
}

func (mock *MockBuilder) Build(args ...string) error {
	return mock.OnBuild(args...)
}

func (mock *MockBuilder) Push(name string) error {
	return mock.OnPush(name)
}
