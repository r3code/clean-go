package engine

import (
	"golang.org/x/net/context"

	"github.com/r3code/clean-go/domain"
)

type (
	// GreetingStorer defines the methods that any
	// data storage provider needs to implement to get
	// and store greetings
	GreetingStorer interface {
		// PutGreeting adds a new Greeting to the datastore
		PutGreeting(c context.Context, greeting *domain.Greeting) error

		// ListGreetings returns existing greetings matching the
		// query provided
		ListGreetings(c context.Context, query *Query) ([]*domain.Greeting, error)
	}
)
