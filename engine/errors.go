package engine

import "errors"

var (
	// ErrGreetingAlredyExists app error
	ErrGreetingAlredyExists = errors.New("Greeting Alredy Exists")
	ErrNotFound             = errors.New("Item Not Found")
	// ErrForbiddenWordPresent app error
	ErrForbiddenWordPresent = errors.New("Greeting has explict content")
	ErrInvalidRequestData   = errors.New("Invalid Request Data")
)
