package engine

import "errors"

var (
	// ErrGreetingAlredyExists app error
	ErrGreetingAlredyExists = errors.New("Greeting Alredy Exists")
	// ErrForbiddenWordPresent app error
	ErrForbiddenWordPresent = errors.New("Greeting has explict content")
)
