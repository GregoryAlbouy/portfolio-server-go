package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store  Store
}

func newServer() *server {
	s := server{
		router: mux.NewRouter().StrictSlash(true),
	}

	s.routes()

	return &s
}

func (s *server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("JSON Marshal error: %v\n", err)
	}

	w.Write(resp)
}
