package engine

import (
	"github.com/r3code/clean-go/domain"
	"golang.org/x/net/context"
)

type (
	// GreetingManager is the interface for our interactor
	GreetingManager interface {
		// Add is the add-a-greeting use-case
		CreateGreeting(c context.Context, r *CreateGreetingRequest) (*CreateGreetingResponse, error)

		// List is the list-the-greetings use-case
		ListGreetings(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error)

		// GetGreeting allows us to read one greeting
		GetGreeting(c context.Context, r *GetGreetingRequest) (*GetGreetingResponse, error)
	}

	// greetingManager implementation
	greetingManager struct {
		repository GreetingStorer
	}
)

// ensure greetingManager realize GreetingManager interface
var _ GreetingManager = &greetingManager{}

type (

	// CreateGreetingRequest to ask engine to add an item
	CreateGreetingRequest struct {
		Author  string
		Content string
	}
	// CreateGreetingResponse from engine after adding the item
	CreateGreetingResponse struct {
		ID int64
	}
	// ListGreetingsRequest ti filer out results
	ListGreetingsRequest struct {
		Offset int
		Count  int
	}
	// ListGreetingsResponse struct with response from the engine
	ListGreetingsResponse struct {
		Greetings []*domain.Greeting
	}
	GetGreetingRequest struct {
		ID int64
	}
	GetGreetingResponse struct {
		Greeting *domain.Greeting
	}
)

func (g *greetingManager) ListGreetings(c context.Context, r *ListGreetingsRequest) (*ListGreetingsResponse, error) {
	q := NewQuery("greetings").Order("date", Descending).Slice(r.Offset, r.Count)
	gl, err := g.repository.ListGreetings(c, q)
	if err != nil {
		return nil, err
	}
	return &ListGreetingsResponse{
		Greetings: gl,
	}, nil
}

func (g *greetingManager) CreateGreeting(c context.Context, r *CreateGreetingRequest) (*CreateGreetingResponse, error) {
	// this is where all our app logic would go - the
	// rules that apply to adding a greeting whether it
	// is being done via the web UI, a console app, or
	// whatever the internet has just been added to ...

	greeting := domain.NewGreeting(r.Author, r.Content)
	err := g.repository.PutGreeting(c, greeting)
	if err != nil {
		return nil, err
	}
	return &CreateGreetingResponse{
		ID: greeting.ID,
	}, nil
}

func (g *greetingManager) GetGreeting(c context.Context, r *GetGreetingRequest) (*GetGreetingResponse, error) {
	q := NewQuery("greetings").Filter("ID", Equal, r.ID)
	gl, err := g.repository.ListGreetings(c, q)
	if err != nil {
		return nil, err
	}
	if len(gl) == 0 {
		return nil, ErrNotFound
	}
	return &GetGreetingResponse{
		Greeting: gl[0],
	}, nil
}
