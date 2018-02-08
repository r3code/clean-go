package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/r3code/clean-go/engine"
)

type (
	mongoDBStorageProvider struct {
		session *mgo.Session
	}
)

// NewStorageProvider creates a new instance of this mongodb storage factory
func NewStorageProvider(url string) (engine.StorageProvider, error) {
	session, err := mgo.DialWithTimeout(url, 10*time.Second)
	if err != nil {
		return nil, err
	}
	err2 := ensureIndexes(session)
	if err2 != nil {
		return nil, err2
	}
	return &mongoDBStorageProvider{session}, nil
}

// CloseStorage closes session
func (sp *mongoDBStorageProvider) CloseStorage() error {
	if sp.session != nil {
		sp.session.Close()
	}
	return nil
}

// NewGreetingRepository creates a new datastore greeting repository
func (sp *mongoDBStorageProvider) NewGreetingRepository() engine.GreetingStorer {
	return newGreetingRepository(sp.session)
}

func ensureIndexes(s *mgo.Session) error {
	index := mgo.Index{
		Key:        []string{"date"},
		Background: true,
	}
	c := s.DB("").C(greetingCollection)
	return c.EnsureIndex(index)
}
