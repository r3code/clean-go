package engine

import (
	"github.com/r3code/clean-go/domain"
	"golang.org/x/net/context"
)

type (
	// GreetingManager is the interface for our interactor
	GreetingManager interface {
		// Add is the add-a-greeting use-case
		CreateGreeting(c context.Context, r *AddGreetingRequest) (*AddGreetingResponse, error)

		// List is the list-the-greetings use-case
		ListGreetings(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error)
	}

	// greetingManager implementation
	greetingManager struct {
		repository GreetingStorer
	}
)

// ensure greetingManager realize GreetingManager interface
var _ GreetingManager = &greetingManager{}

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

func (g *greetingManager) ListGreetings(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error) {
	q := NewQuery("greeting").Order("date", Descending).Slice(0, r.Count)
	gl, err := g.repository.ListGreetings(c, q)
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

func (g *greetingManager) CreateGreeting(c context.Context, r *AddGreetingRequest) (*AddGreetingResponse, error) {
	// this is where all our app logic would go - the
	// rules that apply to adding a greeting whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...

	greeting := domain.NewGreeting(r.Author, r.Content)
	err := g.repository.PutGreeting(c, greeting)
	if err != nil {
		return nil, err
	}
	return &AddGreetingResponse{
		ID: greeting.ID,
	}, nil
}
