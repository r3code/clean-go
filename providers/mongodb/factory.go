package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/r3code/clean-go/engine"
)

type (
	storageFactory struct {
		session *mgo.Session
	}
)

// NewStorage creates a new instance of this mongodb storage factory
func NewStorage(url string) (engine.StorageFactory, error) {
	session, err := mgo.DialWithTimeout(url, 10*time.Second)
	if err != nil {
		return nil, err
	}
	err2 := ensureIndexes(session)
	if err2 != nil {
		return nil, err2
	}
	return &storageFactory{session}, nil
}

// CloseStorage closes session
func (f *storageFactory) CloseStorage() error {
	if f.session != nil {
		f.session.Close()
	}
	return nil
}

// NewGreetingRepository creates a new datastore greeting repository
func (f *storageFactory) NewGreetingRepository() engine.GreetingRepository {
	return newGreetingRepository(f.session)
}

func ensureIndexes(s *mgo.Session) error {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("").C(greetingCollection)
	return c.EnsureIndex(index)
}
