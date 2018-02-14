package engine

type (
	// engine factory stores the state of our engine
	// which only involves a storage factory instance
	// realizes ServiceCreator interface
	serviceEngine struct {
		storageProvider StorageProvider
	}
)

// NewEngine creates a new engine factory that will
// make use of the passed in StorageProvider for any
// data persistence needs.
func NewEngine(sp StorageProvider) ServiceCreator {
	return &serviceEngine{sp}
}

// NewGreetingManager creates a new Greeter interactor wired up
// to use the greetingManager repository from the storage provider
// that the engine has been setup to use.
// @implements ServiceCreator
func (se *serviceEngine) NewGreetingManager() GreetingManager {
	return &greetingManager{
		repository: se.storageProvider.NewGreetingRepository(),
	}
}
