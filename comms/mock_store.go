package comms

type MockStore struct {
	OnStore             func(string, ImageMetadata) error
	OnList              func(string) (*ImageList, error)
	OnCleanup           func()
	OnDelete            func(string) (int, error)
	OnGet               func(string) (*Image, error)
	OnGetLatest         func(string) (*Image, error)
	OnStartTransaction  func() error
	OnCommitTransaction func() error
}

func (mock *MockStore) Store(imageName string, meta ImageMetadata) error {
	return mock.OnStore(imageName, meta)
}

func (mock *MockStore) List(imageName string) (*ImageList, error) {
	return mock.OnList(imageName)
}

func (mock *MockStore) Cleanup() {
	mock.OnCleanup()
}

func (mock *MockStore) Delete(imageName string) (int, error) {
	return mock.OnDelete(imageName)
}

func (mock *MockStore) Get(imageName string) (*Image, error) {
	return mock.OnGet(imageName)
}

func (mock *MockStore) GetLatest(imageName string) (*Image, error) {
	return mock.OnGetLatest(imageName)
}

func (mock *MockStore) StartTransaction() error {
	return mock.OnStartTransaction()
}

func (mock *MockStore) CommitTransaction() error {
	return mock.OnCommitTransaction()
}
