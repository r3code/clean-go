package boltdb

import (
	"encoding/json"
	"strings"

	"github.com/boltdb/bolt"
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

func newGreetingRepository(session *bolt.DB) engine.GreetingStorer {
	return &greetingRepository{session}
}

func (r *greetingRepository) PutGreeting(ctx context.Context, g *domain.Greeting) error {
	// Start read-write transaction.
	tx, err := r.session.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify greeting doesn't already exist.
	b := tx.Bucket([]byte(greetingCollection))
	if v := b.Get(i64tob(g.ID)); v != nil {
		return engine.ErrGreetingAlredyExists
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

func (r *greetingRepository) ListGreetings(ctx context.Context, query *engine.Query) ([]*domain.Greeting, error) {
	var condition func(actual, expected int64) bool = nil
	var idFilter *engine.FilterCondition = nil
	for _, filter := range query.Filters {
		if strings.ToUpper(filter.Property) == "ID" {
			idFilter = filter
			switch filter.Condition {
			case engine.Equal:
				condition = func(actual, expected int64) bool {
					return actual == expected
				}
			}
		}
	}

	gList := []*domain.Greeting{}
	err := r.session.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(greetingCollection))

		cr := b.Cursor()

		for key, value := cr.First(); key != nil; key, value = cr.Next() {
			if condition != nil {
				t := condition(int64(idFilter.Value.(int64)), btoi64(key))
				// fmt.Printf("v = %d, key = %d \n", int64(idFilter.Value.(int64)), btoi64(key))
				if idFilter != nil && !t {
					continue
				}
			}

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
	// TODO: process order query params
	startI := 0
	endI := len(gList)
	if query.Offset > 0 {
		startI = query.Offset
		if query.Offset > len(gList) {
			startI = len(gList)
		}
	}

	if query.Limit > 0 {
		endI = query.Limit
		if query.Offset > 0 {
			endI += query.Offset
		}
		if endI > len(gList) {
			endI = len(gList)
		}
	}
	return gList[startI:endI], nil
}
