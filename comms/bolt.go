package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bradfitz/slice"
	"github.com/spf13/viper"
)

const (
	ImageBucket = "images"
)

// BoltStorage implements GzrMetadataStore and has an un-exported bolt.db pointer
type BoltStorage struct {
	db        *bolt.DB
	activeTxn *bolt.Tx
}

// NewBoltStorage initializes a BoltDB connection, makes sure the correct buckets exist,
// and returns a BoltStorage pointer with the established connection
func NewBoltStorage() (GzrMetadataStore, error) {
	dbPath := viper.GetString("datastore.db_path")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	store := &BoltStorage{db: db}
	txn, err := store.db.Begin(true)
	if err != nil {
		store.Cleanup()
		return nil, err
	}
	_, err = txn.CreateBucketIfNotExists([]byte(ImageBucket))
	if err != nil {
		store.Cleanup()
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		store.Cleanup()
		return nil, err
	}
	return store, nil
}

// List queries the Bolt store for all images stored under a particular name
func (store *BoltStorage) List(imageName string) (*ImageList, error) {
	var images []*Image
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		c := b.Cursor()
		prefix := []byte(imageName)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			img := store.extractImage(v, k)
			images = append(images, img)
		}
		return nil
	})
	if err != nil {
		return &ImageList{}, err
	}
	return &ImageList{Images: images}, nil
}

// Store stores the metadata about an image associated with its name
func (store *BoltStorage) Store(imageName string, meta ImageMetadata) error {
	b := store.activeTxn.Bucket([]byte(ImageBucket))
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	key, err := store.createKey(imageName)
	if err != nil {
		return err
	}
	err = b.Put(key, data)
	if err != nil {
		return err
	}
	return nil
}

// Cleanup closes the Bolt connection
func (store *BoltStorage) Cleanup() {
	store.db.Close()
}

// Delete deletes all information related to IMAGE_NAME:VERSION
func (store *BoltStorage) Delete(imageName string) (int, error) {
	b := store.activeTxn.Bucket([]byte(ImageBucket))
	c := b.Cursor()
	prefix := []byte(imageName)
	deleted := 0
	for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
		err := b.Delete(k)
		if err != nil {
			return 0, err
		}
		deleted += 1
	}
	return deleted, nil
}

// Get returns a single image based on a name
func (store *BoltStorage) Get(imageName string) (*Image, error) {
	var image *Image
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		c := b.Cursor()
		prefix := []byte(imageName)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			image = store.extractImage(v, k)
		}
		return nil
	})
	if err != nil {
		return &Image{}, err
	}
	return image, nil
}

// GetLatest returns the latest image from a name
func (store *BoltStorage) GetLatest(imageName string) (*Image, error) {
	images, err := store.List(imageName)
	if err != nil {
		return nil, err
	}
	slice.Sort(images.Images, func(i, j int) bool {
		return images.Images[j].Meta.CreatedAt < images.Images[i].Meta.CreatedAt
	})
	return images.Images[0], nil
}

// StartTransaction starts a new Bolt transaction and adds it to the Storage
func (store *BoltStorage) StartTransaction() error {
	bTxn, err := store.db.Begin(true)
	if err != nil {
		return err
	}
	store.activeTxn = bTxn
	return nil
}

// CommitTransaction commits the active transaction
func (store *BoltStorage) CommitTransaction() error {
	return store.activeTxn.Commit()
}

// createKey creates the key used to tag data in Bolt
func (store *BoltStorage) createKey(imageName string) ([]byte, error) {
	splitName := strings.Split(imageName, ":")
	if len(splitName) != 2 {
		return []byte{}, fmt.Errorf("IMAGE_NAME must be formatted as NAME:VERSION and must contain only the seperating colon")
	}
	now := time.Now().Format("20060102")
	return []byte(fmt.Sprintf("%s:%s:%s", splitName[0], splitName[1], now)), nil
}

// extractImage transforms raw []byte of metadata and key into a full Image
func (store *BoltStorage) extractImage(data []byte, key []byte) *Image {
	var meta ImageMetadata
	json.Unmarshal(data, &meta)
	return &Image{Name: string(key), Meta: meta}
}
