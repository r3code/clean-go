package apperror

import (
	"errors"
)

var (
	GreetingAlredyExists = errors.New("Greeting Alredy Exists")
)
