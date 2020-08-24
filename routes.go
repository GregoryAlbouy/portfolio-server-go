package main

import (
	"net/http"
)

// RouteConfig struct
type RouteConfig struct {
	path    string
	method  string
	handler http.HandlerFunc
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.handleProjectRoutes()
	// s.handleContactRoutes()
}

func (s *server) handleProjectRoutes() {
	sub := s.router.PathPrefix("/api/v1/projects").Subrouter()

	projectRoutes := []RouteConfig{
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

// TODO: handlers
func (s *server) handleContactRoutes() {
	sub := s.router.PathPrefix("/contact").Subrouter()

	contactRoutes := []RouteConfig{
		{path: "/", method: "GET", handler: nil},
		{path: "/{id}", method: "POST", handler: nil},
		{path: "/{id}", method: "GET", handler: nil},
		{path: "/{id}", method: "DELETE", handler: nil},
	}

	for _, r := range contactRoutes {
		sub.HandleFunc(r.path, r.handler).Methods(r.method)
	}
}

func (s *server) handleAuthRoutes() {
	ct := s.router.PathPrefix("/contact").Subrouter()
	ct.HandleFunc("/", nil).Methods("GET")
	ct.HandleFunc("/", nil).Methods("POST")
	ct.HandleFunc("/{id}", nil).Methods("GET")
}
