// +build !appengine

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/r3code/clean-go/providers/boltdb"

	"github.com/r3code/clean-go/adapters/web"
	"github.com/r3code/clean-go/engine"
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

	endpoint := ":8080"
	router := web.NewWebAdapter(e, true)
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	const shutdownTimeout = 8 * time.Second
	go gracefulShutdown(server, shutdownTimeout)
	log.Println("Listen to http at " + endpoint)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}

// gracefulShutdown catch interrupt signal and wait server to stop,
// but not longer than `timeout`
func gracefulShutdown(s *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.Shutdown(ctx)
}
