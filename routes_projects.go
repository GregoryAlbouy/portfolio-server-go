package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) handleProjectList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := s.store.GetProjectList()
		if err != nil {
			fmt.Print(err)
			s.respond(w, r, nil, http.StatusNoContent)
		}

		s.respond(w, r, projects.toJSON(), http.StatusOK)
	}
}

func (s *server) handleProjectDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			fmt.Println("Cannot load movie with id: ", vars["id"])
		}
		resp := map[string]interface{}{
			"id":    id,
			"query": "project detail",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}

func (s *server) handleProjectCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"message": "TODO",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}

func (s *server) handleProjectUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"message": "TODO",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}

func (s *server) handleProjectDeletion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"message": "TODO",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}
