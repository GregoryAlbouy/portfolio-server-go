package main

import (
	"fmt"
	"gregoryalbouy-server-go/clog"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = "8080"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	// server
	srv := newServer()
	srv.store = &dbStore{}
	port := port()
	addr := ":" + port

	// store
	if err = srv.store.Open(); err != nil {
		return
	}
	defer srv.store.Close()

	// serve
	if err = http.ListenAndServe(addr, srv.router); err != nil {
		return
	}

	return
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("Listening to port %s\n", clog.Blue(port))
	return port
}
