package mongodb

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/r3code/clean-go/domain"
	"github.com/r3code/clean-go/engine"
)

type (
	greetingRepository struct {
		session *mgo.Session
	}
)

var (
	greetingCollection = "greeting"
)

func newGreetingRepository(session *mgo.Session) engine.GreetingStorer {
	return &greetingRepository{session}
}

func (r greetingRepository) PutGreeting(c context.Context, g *domain.Greeting) error {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("").C(greetingCollection)
	if g.ID == 0 {
		g.ID = getNextSequence(s, greetingCollection)
	}
	if _, err := col.Upsert(bson.M{"_id": g.ID}, g); err != nil {
		return err
	}
	return nil
}

func (r greetingRepository) ListGreetings(c context.Context, query *engine.Query) ([]*domain.Greeting, error) {
	s := r.session.Clone()
	defer s.Close()

	col := s.DB("").C(greetingCollection)
	g := []*domain.Greeting{}
	q := translateQuery(col, query)
	if err := q.All(&g); err != nil {
		return nil, err
	}
	return g, nil
}
