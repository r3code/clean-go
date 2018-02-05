package engine

import (
	"github.com/r3code/clean-go/domain"
	"golang.org/x/net/context"
)

type (
	// Greeter is the interface for our interactor
	Greeter interface {
		// Add is the add-a-greeting use-case
		Add(c context.Context, r *AddGreetingRequest) (*AddGreetingResponse, error)

		// List is the list-the-greetings use-case
		List(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error)
	}

	// greeter implementation
	greeter struct {
		repository GreetingRepository
	}
)

// ensure greeter realize Greeter interface
var _ Greeter = &greeter{}

// NewGreeter creates a new Greeter interactor wired up
// to use the greeter repository from the storage provider
// that the engine has been setup to use.
func (f *engineServiceFactory) NewGreeter() Greeter {
	return &greeter{
		repository: f.NewGreetingRepository(),
	}
}

type (
	// ListGreetingsRequest ti filer out results
	ListGreetingsRequest struct {
		Count int
	}
	// ListGreetingsResponse struct with response from the engine
	ListGreetingsResponse struct {
		Greetings []*domain.Greeting
	}
)

func (g *greeter) List(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error) {
	q := NewQuery("greeting").Order("date", Descending).Slice(0, r.Count)
	gl, err := g.repository.List(c, q)
	if err != nil {
		return nil, err
	}
	return &ListGreetingsResponse{
		Greetings: gl,
	}, nil
}

type (
	// AddGreetingRequest to ask engine to add an item
	AddGreetingRequest struct {
		Author  string
		Content string
	}
	// AddGreetingResponse from engine after adding the item
	AddGreetingResponse struct {
		ID int64
	}
)

func (g *greeter) Add(c context.Context, r *AddGreetingRequest) (*AddGreetingResponse, error) {
	// this is where all our app logic would go - the
	// rules that apply to adding a greeting whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...

	greeting := domain.NewGreeting(r.Author, r.Content)
	err := g.repository.Put(c, greeting)
	if err != nil {
		return nil, err
	}
	return &AddGreetingResponse{
		ID: greeting.ID,
	}, nil
}
