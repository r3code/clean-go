// +build !appengine

package main

import (
	"log"
	"net/http"

	"github.com/r3code/clean-go/adapters/web"
	"github.com/r3code/clean-go/engine"
	"github.com/r3code/clean-go/providers/mongodb"
)

// when running in traditional or 'standalone' mode
// we're going to use MongoDB as the storage provider
// and start the webserver running ourselves.
func main() {
	s, err := mongodb.NewStorageProvider(config.MongoURL)
	if err != nil {
		log.Fatalln("Storage init error: " + err.Error())
	}
	defer func() {
		cerr := st.CloseStorage()
		if cerr != nil {
			log.Fatalln("Storage close error: " + cerr.Error())
		}
	}()
	e := engine.NewEngine(s)
	http.ListenAndServe(":8080", web.NewWebAdapter(e, true))
}
