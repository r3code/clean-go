package cli

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"golang.org/x/net/context"

	"github.com/r3code/clean-go/engine"
	"github.com/urfave/cli"
)

type (
	greetingManager struct {
		engine.GreetingManager
	}
)

// wire up the greetings routes
func initGreetings(a *cli.App, f engine.ServiceCreator) {
	greetingManager := &greetingManager{f.NewGreetingManager()}
	a.Commands = append(a.Commands, []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list saved greetings",
			Action:  greetingManager.listGreetings,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "<author> <text> add a greeting to the list,",
			Action:  greetingManager.createGreeting,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "author, a",
					Usage: "author `NAME`",
				},
				cli.StringFlag{
					Name:  "content, c",
					Usage: "greeting `CONTENT`",
				},
			},
		},
	}...)

}

func (g *greetingManager) createGreeting(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New(`You must pass two args "<author_name>" "<text>"`)
	}
	req := &engine.AddGreetingRequest{
		Author:  c.Args().Get(0),
		Content: c.Args().Get(1),
	}
	// We can set auth params here for example and pass them in the context
	// ctx := context.WithValue(context.Background(), "userid", c.String("user") )
	// ctx := context.WithValue(context.Background(), "pass", c.String("pass") )
	ctx := context.Background()
	res, err := g.CreateGreeting(ctx, req)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("added new greeting ID=%d\n", res.ID)
	return nil
}

func (g *greetingManager) listGreetings(c *cli.Context) error {
	ctx := context.Background()
	count, err := strconv.Atoi(c.Args().Get(0))
	if err != nil || count == 0 {
		count = 5
	}
	req := &engine.ListGreetingsRequest{
		Count: count,
	}
	res, err := g.ListGreetings(ctx, req)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	log.Printf("List of greetings: \n")
	for i := range res.Greetings {
		greet := res.Greetings[i]
		log.Printf("ID: %d; Date: %s, Author: %s; Text: %s\n", greet.ID, greet.Date, greet.Author, greet.Content)
	}
	return nil
}
