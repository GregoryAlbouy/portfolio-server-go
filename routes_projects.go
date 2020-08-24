package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func (s *server) readProjectList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := s.store.GetProjectList()
		if err != nil {
			fmt.Print(err)
			s.respond(w, r, nil, http.StatusNoContent)
		}

		s.respond(w, r, projects.formatJSON(), http.StatusOK)
	}
}

func (s *server) readProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["id"]
		resp, err := s.store.GetProjectBySlug(slug)
		if err != nil {
			s.respond(w, r, nil, http.StatusNotFound)
			fmt.Println(err)
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

// func (s *server) readProjectByID() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		id, err := strconv.ParseInt(vars["id"], 10, 64)
// 		if err != nil {
// 			fmt.Println("Cannot load movie with id: ", vars["id"])
// 		}
// 		resp := map[string]interface{}{
// 			"id":    id,
// 			"query": "project detail",
// 		}
// 		s.respond(w, r, resp, http.StatusNotImplemented)
// 	}
// }

func (s *server) createProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := Project{}

		if err := r.ParseMultipartForm(0); err != nil {
			fmt.Printf("PARSE FORM ERROR: %+v\n", err)
			s.respond(w, r, "Invalid input", http.StatusBadRequest)
		}

		if err := decoder.Decode(&p, r.PostForm); err != nil {
			fmt.Printf("DECODE ERROR (PostForm): %+v\n", err)
			s.respond(w, r, "Invalid input", http.StatusBadRequest)
		}

		fmt.Printf("%+v", p)
		// p.toDB()
		// p.toJSON()

		// if err := s.decode(w, r, &p); err != nil {
		// 	s.respond(w, r, "Invalid request", http.StatusBadRequest)
		// 	log.Panicln("DECODE ISSUE", err)
		// 	return
		// }

		if !p.IsValid() {
			fmt.Println("Project is invalid")
			s.respond(w, r, "Project data is incomplete", http.StatusBadRequest)
		}

		if err := s.store.CreateProject(&p); err != nil {
			s.respond(w, r, "Invalid request", http.StatusBadRequest)
			fmt.Println("CREATE ISSUE", err)
			return
		}

		s.respond(w, r, p, http.StatusCreated)
	}
}

func (s *server) updateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.store.UpdateProjectBySlug("codemon")

		resp := map[string]interface{}{
			"message": "TODO",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}

func (s *server) deleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"message": "TODO",
		}
		s.respond(w, r, resp, http.StatusNotImplemented)
	}
}
