package main

import (
	"fmt"
	"gregoryalbouy-server-go/clog"
	"log"
	"net/http"
	"os"
	"time"
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
	start := time.Now()
	// display env
	fmt.Println("Environment", clog.Blue(os.Getenv("APP_ENV")))

	// server
	srv := newServer().attachStore(&dbStore{})
	port := port()
	addr := ":" + port

	// store
	if err = srv.store.Open(); err != nil {
		return
	}
	defer srv.store.Close()

	fmt.Printf("%s (%s) http://127.0.0.1%s \n", clog.Green("Ready"), time.Since(start), addr)

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
	fmt.Println("Port", clog.Blue(port))

	return port
}
