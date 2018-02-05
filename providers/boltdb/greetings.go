package boltdb

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/r3code/clean-go/apperror"
	"golang.org/x/net/context"

	"github.com/r3code/clean-go/domain"
	"github.com/r3code/clean-go/engine"
)

type (
	greetingRepository struct {
		session *bolt.DB
	}
)

var (
	greetingCollection = "greeting"
)

func newGreetingRepository(session *bolt.DB) engine.GreetingRepository {
	return &greetingRepository{session}
}

func (r greetingRepository) Put(ctx context.Context, g *domain.Greeting) error {
	// Start read-write transaction.
	tx, err := r.session.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify greeting doesn't already exist.
	b := tx.Bucket([]byte(greetingCollection))
	if v := b.Get(i64tob(g.ID)); v != nil {
		return apperror.GreetingAlredyExists
	}

	id, _ := b.NextSequence()
	g.ID = int64(id)

	// Marshal and insert record.
	if v, err := json.Marshal(g); err != nil {
		return err
	} else if err := b.Put(i64tob(g.ID), v); err != nil {
		return err
	}

	return tx.Commit()
}

func (r greetingRepository) List(ctx context.Context, query *engine.Query) ([]*domain.Greeting, error) {
	gList := []*domain.Greeting{}
	err := r.session.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(greetingCollection))

		cr := b.Cursor()

		for key, value := cr.First(); key != nil; key, value = cr.Next() {
			var g domain.Greeting
			if err := json.Unmarshal(value, &g); err != nil {
				return err
			}
			gList = append(gList, &g)
		}

		return nil // from func(tx *bolt.Tx)
	})
	if err != nil {
		return nil, err
	}
	return gList, nil
}
