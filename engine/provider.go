package engine

type (
	// StorageProvider is the interface that a storage
	// provider needs to implement so that the engine can
	// request repository instances as it needs them
	StorageProvider interface {
		// NewGreetingStorer returns a storage specific
		// GreetingStorer implementation
		NewGreetingRepository() GreetingStorer
		CloseStorage() error
	}
)
