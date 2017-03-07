package comms

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
)

// GzrMetadataStore is the interface that should be implemented for any
// backend that needs to handle image data storage
type GzrMetadataStore interface {
	// Store stores image metadata with a name
	Store(string, ImageMetadata) error
	// List lists all of the images under a name
	List(string) (*ImageList, error)
	// Cleanup allows the storage backend to clean up any connections, etc
	Cleanup()
	// Delete deletes all images under a nmae
	Delete(string) error
	// Get gets a single image with a version
	Get(string) (*Image, error)
	// NewTransaction returns a new StorageTransaction
	NewTransaction() (StorageTransaction, error)
}

type StorageTransaction interface {
	Commit() error
}

// Image is a struct unifying an image name with its metadata
type Image struct {
	// Name is the image's full name
	Name string `json:"name"`
	// Meta is the metadata related to the image
	Meta ImageMetadata `json:"metadata"`
}

// ImageMetadata is a struct containing necessary metadata about a particular image
type ImageMetadata struct {
	// GitCommit is the commit related to the built image
	GitCommit string `json:"git-commit"`
	// GitTag is the tag related to the built image if it exists
	GitTag []string `json:"git-tag"`
	// GitAnnotation is the annotation related to the built image if it exists
	GitAnnotation []string `json:"git-annotation"`
	// GitOrigin is the remote origin for the git repo associated with the image
	GitOrigin string `json:"git-origin"`
	// CreatedAt is the time the metadata was stored, with day granularity
	CreatedAt string `json:"created-at"`
}

// ImageList is a collection of Images
type ImageList struct {
	Images []*Image `json:"images"`
}

// SerializeForCLI takes an io.Writer and writes templatized data to it representing an image list
func (l *ImageList) SerializeForCLI(wr io.Writer) error {
	return l.cliTemplate().Execute(wr, l)
}

func (l *ImageList) cliTemplate() *template.Template {
	t := template.New("Images")
	t, _ = t.Parse(`Images {{range .Images}}
- name: {{.Name}}
  -- git-commit: {{.Meta.GitCommit}}
  -- git-tag: [{{ range $index, $element := .Meta.GitTag}}{{if $index}}, {{end}}{{$element}}{{end}}]
  -- git-annotation: [{{ range $index, $element := .Meta.GitAnnotation}}{{if $index}}, {{end}}{{$element}}{{end}}]
  -- git-origin: {{.Meta.GitOrigin}}
  -- created-at: {{.Meta.CreatedAt}}
{{end}}
`)
	return t
}

// SerializeForWire returns a JSON representation of the ImageList
func (l *ImageList) SerializeForWire() ([]byte, error) {
	return json.Marshal(l)
}

// SerializeForWire returns a JSON representation of the Image
func (i *Image) SerializeForWire() ([]byte, error) {
	return json.Marshal(i)
}

// CreateMeta takes a ReadWriter and returns an instance of ImageMetadata
// after parsing
func CreateMeta(reader io.ReadWriter) (ImageMetadata, error) {
	var meta ImageMetadata
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("Could not read metadata file: %s", err.Error())
	}

	err = json.Unmarshal(b, &meta)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("Could not read metadata file: %s\nCheck the data types in your image metadata JSON.", err.Error())
	}
	return meta, nil
}
