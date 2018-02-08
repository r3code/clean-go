// +build appengine

package appengine

import (
	"github.com/r3code/clean-go/engine"
)

type (
	storageFactory struct{}
)

// NewStorageProvider creates a new instance of this datastore storage factory
func NewStorageProvider() (engine.StorageProvider, error) {
	return &storageFactory{}, nil
}

// NewGreetingRepository creates a new datastore greeting repository
func (f *storageFactory) NewGreetingRepository() engine.GreetingStorer {
	return newGreetingRepository()
}

// CloseStorage closes session
func (f *storageFactory) CloseStorage() error {
	return nil
}
