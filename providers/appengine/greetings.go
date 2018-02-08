// +build appengine

package appengine

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/r3code/clean-go/domain"
	"github.com/r3code/clean-go/engine"
)

type (
	greetingRepository struct{}

	// greeting is the internal struct we use for persistence
	// it allows us to attach the datastore PropertyLoadSaver
	// methods to the struct that we don't own
	greeting struct {
		domain.Greeting
	}
)

var (
	greetingKind = "greeting"
)

// var _ engine.GreetingStorer = &greetingRepository{}

func newGreetingRepository() engine.GreetingStorer {
	return &greetingRepository{}
}

func (r greetingRepository) PutGreeting(c context.Context, g *domain.Greeting) error {
	d := &greeting{*g}
	k := datastore.NewIncompleteKey(c, greetingKind, guestbookEntityKey(c))
	_, err := datastore.Put(c, k, d)
	return err
}

func (r greetingRepository) ListGreetings(c context.Context, query *engine.Query) ([]*domain.Greeting, error) {
	g := []*greeting{}
	q := translateQuery(greetingKind, query)
	q = q.Ancestor(guestbookEntityKey(c))

	k, err := q.GetAll(c, &g)
	if err != nil {
		return nil, err
	}
	o := make([]*domain.Greeting, len(g))
	for i := range g {
		o[i] = &g[i].Greeting
		o[i].ID = k[i].IntID()
	}
	return o, nil
}

func guestbookEntityKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "guestbook", "", 1, nil)
}

func (x *greeting) Load(props []datastore.Property) error {
	for _, prop := range props {
		switch prop.Name {
		case "author":
			x.Author = prop.Value.(string)
		case "content":
			x.Content = prop.Value.(string)
		case "date":
			x.Date = prop.Value.(time.Time)
		}
	}
	return nil
}

func (x *greeting) Save() ([]datastore.Property, error) {
	ps := []datastore.Property{
		datastore.Property{Name: "author", Value: x.Author, NoIndex: true, Multiple: false},
		datastore.Property{Name: "content", Value: x.Content, NoIndex: true, Multiple: false},
		datastore.Property{Name: "date", Value: x.Date, NoIndex: false, Multiple: false},
	}
	return ps, nil
}
