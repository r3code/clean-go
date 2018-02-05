package boltdb

import (
	"time"

	"github.com/boltdb/bolt"

	"github.com/r3code/clean-go/engine"
)

type (
	storageFactory struct {
		session *bolt.DB
	}
)

// NewStorage creates a new instance of this mongodb storage factory
func NewStorage(filename string) (engine.StorageFactory, error) {
	session, err := bolt.Open(filename, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err2 := initBuckets(session)
	if err2 != nil {
		return nil, err2
	}
	return &storageFactory{session}, nil
}

// NewGreetingRepository creates a new datastore greeting repository
func (f *storageFactory) NewGreetingRepository() engine.GreetingRepository {
	return newGreetingRepository(f.session)
}

// CloseStorage closes session
func (f *storageFactory) CloseStorage() error {
	if f.session != nil {
		return f.session.Close()
	}
	return nil
}

func initBuckets(session *bolt.DB) error {
	// Initialize top-level buckets.
	tx, err := session.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte(greetingCollection)); err != nil {
		return err
	}

	return tx.Commit()
}
