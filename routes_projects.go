package main

import (
	"fmt"
	"net/http"
)

func (s *server) handleProjectList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := s.store.GetProjects()
		if err != nil {
			fmt.Print(err)
			s.respond(w, r, nil, http.StatusNoContent)
		}

		s.respond(w, r, projects, http.StatusOK)
	}
}
