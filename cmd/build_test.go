package cmd

import (
	"testing"

	"github.com/bypasslane/gzr/comms"
)

var (
	buildCalled bool
	pushCalled  bool
	storeCalled bool
)

func TestBuildImage(t *testing.T) {
	buildCalled = false
	pushCalled = false
	imageStore = &comms.MockStore{
		OnStore: callStore,
	}
	manager := &comms.MockManager{
		OnBuild: callBuild,
		OnPush:  callPush,
	}
	err := buildHandler([]string{}, manager)
	if err != nil {
		t.Errorf("buildImage errored with %s", err.Error())
	}
}

func callBuild(args ...string) error {
	buildCalled = true
	return nil
}

func callPush(name string) error {
	pushCalled = true
	return nil
}

func callStore(name string, meta comms.ImageMetadata) error {
	storeCalled = true
	return nil
}
