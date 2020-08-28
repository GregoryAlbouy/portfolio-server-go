package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type routeConfig struct {
	path    string
	method  string
	handler http.HandlerFunc
}

func (s *server) routes() {
	routes := []routeConfig{
		{path: "/", method: "GET", handler: s.handleIndex()},
		{path: "/token/", method: "POST", handler: s.createToken()},
	}
	s.serveRoutes(s.router, routes)

	s.projectRoutes()
	s.userRoutes()
	s.contactRoutes()
}

func (s *server) projectRoutes() {
	subrouter := s.router.PathPrefix("/api/v1/projects").Subrouter()

	routes := []routeConfig{
		{path: "/", method: "GET", handler: s.getProjectList()},
		{path: "/", method: "POST", handler: s.authOnly(s.createProject())},
		{path: "/{id}", method: "GET", handler: s.getProject()},
		{path: "/{id}", method: "PUT", handler: s.authOnly(s.updateProject())},
		{path: "/{id}", method: "DELETE", handler: s.authOnly(s.deleteProject())},
	}

	s.serveRoutes(subrouter, routes)
}

func (s *server) userRoutes() {
	subrouter := s.router.PathPrefix("/users").Subrouter()
	subrouter.Use(s.adminOnlyMiddleware)

	routes := []routeConfig{
		{path: "/", method: "POST", handler: s.createUser()},
		{path: "/", method: "GET", handler: s.getUserList()},
		{path: "/{id}", method: "GET", handler: s.getUser()},
		{path: "/{id}", method: "DELETE", handler: s.deleteUser()},
	}

	s.serveRoutes(subrouter, routes)
}

// TODO: handlers
func (s *server) contactRoutes() {
	subrouter := s.router.PathPrefix("/contact").Subrouter()

	routes := []routeConfig{
		{path: "/", method: "GET", handler: s.authOnly(s.getMessageList())},
		{path: "/", method: "POST", handler: s.postMessage()},
		{path: "/{id}", method: "DELETE", handler: s.authOnly(s.deleteMessage())},
	}

	s.serveRoutes(subrouter, routes)
}

func (s *server) serveRoutes(router *mux.Router, routes []routeConfig) {
	for _, r := range routes {
		router.HandleFunc(r.path, r.handler).Methods("OPTIONS", r.method)
	}
}
