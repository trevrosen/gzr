package comms

type MockStore struct {
	OnStore   func(string, ImageMetadata) error
	OnList    func(string) (ImageList, error)
	OnCleanup func()
	OnDelete  func(string) error
	OnGet     func(string) (Image, error)
}

func (mock *MockStore) Store(imageName string, meta ImageMetadata) error {
	return mock.OnStore(imageName, meta)
}

func (mock *MockStore) List(imageName string) (ImageList, error) {
	return mock.OnList(imageName)
}

func (mock *MockStore) Cleanup() {
	mock.OnCleanup()
}

func (mock *MockStore) Delete(imageName string) error {
	return mock.OnDelete(imageName)
}

func (mock *MockStore) Get(imageName string) (Image, error) {
	return mock.OnGet(imageName)
}
