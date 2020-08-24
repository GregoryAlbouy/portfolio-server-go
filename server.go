package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store  Store
}

func newServer() *server {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	))

	s := server{
		router: r,
	}

	s.routes()

	return &s
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

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
