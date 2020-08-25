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
	s.handleAuthRoutes()
	// s.handleContactRoutes()
}

func (s *server) handleProjectRoutes() {
	sub := s.router.PathPrefix("/api/v1/projects").Subrouter()

	projectRoutes := []RouteConfig{
		{path: "/", method: "GET", handler: s.getProjectList()},
		{path: "/", method: "POST", handler: s.auth(s.createProject())},
		{path: "/{id}", method: "GET", handler: s.getProject()},
		{path: "/{id}", method: "PUT", handler: s.auth(s.updateProject())},
		{path: "/{id}", method: "DELETE", handler: s.auth(s.deleteProject())},
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
	sub := s.router.PathPrefix("/auth").Subrouter()

	authRoutes := []RouteConfig{
		{path: "/", method: "POST", handler: s.createToken()},
	}

	for _, r := range authRoutes {
		sub.HandleFunc(r.path, r.handler).Methods(r.method)
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}
}
