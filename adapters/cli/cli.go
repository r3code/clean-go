package cli

import (
	"github.com/r3code/clean-go/engine"
	"github.com/urfave/cli"
)

// NewCliAdapter a new web adaptor which will
// handle the command line interface and pass calls on to the
// engine to do the real work (that's why the engine
// factory is passed in - so anything that *it* needs
// is unknown to this).
// Because the command line adapter ends up quite lightweight
// it easier to replace. Changing should be far less costly.
func NewCliAdapter(f engine.ServiceFactory, log bool) *cli.App {
	app := cli.NewApp()
	app.Usage = "greeting book"
	app.Commands = []cli.Command{}
	app.Before = beforeAction
	initGreetings(app, f)
	return app
}

func beforeAction(c *cli.Context) error { return nil }
