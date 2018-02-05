// +build appengine

package appengine

import (
	"github.com/r3code/clean-go/engine"
)

type (
	storageFactory struct{}
)

// NewStorage creates a new instance of this datastore storage factory
func NewStorage() (engine.StorageFactory, error) {
	return &storageFactory{}, nil
}

// NewGreetingRepository creates a new datastore greeting repository
func (f *storageFactory) NewGreetingRepository() engine.GreetingRepository {
	return newGreetingRepository()
}

// CloseStorage closes session
func (f *storageFactory) CloseStorage() error {
	return nil
}
