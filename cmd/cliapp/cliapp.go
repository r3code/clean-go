package main

import (
	"fmt"
	"log"
	"os"

	"github.com/r3code/clean-go/adapters/cli"
	"github.com/r3code/clean-go/engine"
	"github.com/r3code/clean-go/providers/boltdb"
)

// when running in traditional or 'standalone' mode
// we're going to use BoltDB as the storage provider
// and start the webserver running ourselves.
func main() {
	st, err := boltdb.NewStorage(config.BoltDBFile)
	if err != nil {
		log.Fatalln("Storage init error: " + err.Error())
	}
	defer st.CloseStorage()

	e := engine.NewEngine(st)
	app := cli.NewCliAdapter(e, true)
	app.Name = "greeting-cli"
	app.Version = "0.1"

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
