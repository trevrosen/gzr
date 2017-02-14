package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

const (
	ImageBucket = "images"
)

// BoltStorage implements GozerMetadataStore and has an exported bolt.DB pointer
type BoltStorage struct {
	DB *bolt.DB
}

// NewBoltStorage initializes a BoltDB connection, makes sure the correct buckets exist,
// and returns a BoltStorage pointer with the established connection
func NewBoltStorage() (GozerMetadataStore, error) {
	store := &BoltStorage{}
	dbPath := viper.GetString("datastore.db_path")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	store.DB = db
	txn, err := store.DB.Begin(true)
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
func (store *BoltStorage) List(imageName string) ([]Image, error) {
	var images []Image
	err := store.DB.View(func(tx *bolt.Tx) error {
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
		return []Image{}, err
	}
	return images, nil
}

// Store stores the metadata about an image associated with its name
func (store *BoltStorage) Store(imageName string, meta ImageMetadata) error {
	return store.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
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
	})
}

// Cleanup closes the Bolt connection
func (store *BoltStorage) Cleanup() {
	store.DB.Close()
}

// Delete deletes all information related to IMAGE_NAME:VERSION
func (store *BoltStorage) Delete(imageName string) error {
	return store.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		c := b.Cursor()
		prefix := []byte(imageName)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			err := b.Delete(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Get returns a single image based on a name
func (store *BoltStorage) Get(imageName string) (Image, error) {
	var image Image
	err := store.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ImageBucket))
		c := b.Cursor()
		prefix := []byte(imageName)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			image = store.extractImage(v, k)
		}
		return nil
	})
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

// createKey creates the key used to tag data in Bolt
func (store *BoltStorage) createKey(imageName string) ([]byte, error) {
	splitName := strings.Split(imageName, ":")
	if len(splitName) != 2 {
		return []byte{}, fmt.Errorf("IMAGE_NAME must be formatted as NAME:VERSION and must contain only the seperating colon")
	}
	now := time.Now()
	nowString := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	return []byte(fmt.Sprintf("%s:%s:%s", splitName[0], splitName[1], nowString)), nil
}

// extractImage transforms raw []byte of metadata and key into a full Image
func (store *BoltStorage) extractImage(data []byte, key []byte) Image {
	var meta ImageMetadata
	json.Unmarshal(data, &meta)
	return Image{Name: string(key), Meta: meta}
}
