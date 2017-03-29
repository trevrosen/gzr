package comms

import "testing"

const (
	annotatedTags = `
v0.0.0    my annotation
v0.1.0    my 2nd annotation`

	notAnnotatedTags = `
v0.0.0
v0.1.0`
)

func TestProcessTags_Annotated_OK(t *testing.T) {
	tags, annotations := processTags(annotatedTags)

	if len(tags) != 2 {
		t.Errorf("Expected to extract 2 tags, but extracted %d", len(tags))
	}

	if len(annotations) != 2 {
		t.Errorf("Expected to extract 2 annotations, but extracted %d", len(annotations))
	}
}

func TestProcessTags_NotAnnotated_OK(t *testing.T) {
	tags, annotations := processTags(notAnnotatedTags)

	if len(tags) != 2 {
		t.Errorf("Expected to extract 2 tags, but extracted %d", len(tags))
	}

	if len(annotations) != 0 {
		t.Errorf("Expected to extract 0 annotations, but extracted %d", len(annotations))
	}
}
