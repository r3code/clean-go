package engine

type (
	// ServiceCreator interface allows us to provide
	// other parts of the system with a way to make
	// instances of our use-case / interactors when
	// they need to
	ServiceCreator interface {
		// NewGreetingManager creates a new GreetingManager interactor
		NewGreetingManager() GreetingManager
	}
)
