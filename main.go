package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gregoryalbouy/portfolio-server-go/clog"
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
	// display env
	fmt.Println("Environment", clog.Blue(os.Getenv("APP_ENV")))

	// server
	s := newServer().
		attachStore(&dbStore{}).
		setPortWithDefault(defaultPort)

	// store
	if err = s.store.Open(); err != nil {
		return
	}
	defer s.store.Close()

	// serve
	if err = s.serve(); err != nil {
		return
	}

	return
}
