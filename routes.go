package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ServerRoute struct
type ServerRoute struct {
	path    string
	handler http.HandlerFunc
	method  string
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.projectRoutes()
}

func (s *server) projectRoutes() {
	sub := s.router.PathPrefix("/api/v1/projects").Subrouter()
	projectRoutes := []ServerRoute{
		{path: "/", method: "GET", handler: s.readProjectList()},
		{path: "/", method: "POST", handler: s.createProject()},
		{path: "/{id}", method: "GET", handler: s.readProject()},
		{path: "/{id}", method: "PUT", handler: s.updateProject()},
		{path: "/{id}", method: "DELETE", handler: s.deleteProject()},
	}

	for _, r := range projectRoutes {
		sub.HandleFunc(r.path, r.handler).Methods(r.method)
	}
}
func (s *server) contactRoutes() {

}

func (s *server) authRoutes() {
	ct := s.router.PathPrefix("/contact").Subrouter()
	ct.HandleFunc("/", nil).Methods("GET")
	ct.HandleFunc("/", nil).Methods("POST")
	ct.HandleFunc("/{id}", nil).Methods("GET")
}

func (s *server) serveRoutes(router *mux.Router, routes []ServerRoute) {
}
