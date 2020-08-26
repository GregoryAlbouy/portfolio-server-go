package main

import (
	"fmt"
	"gregoryalbouy-server-go/clog"
	"gregoryalbouy-server-go/utl"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) getProjectList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := s.store.GetProjectList()
		if err != nil {
			s.respond(w, r, "No result", http.StatusNoContent)
			return
		}

		s.respond(w, r, projects.formatJSON(), http.StatusOK)
	}
}

func (s *server) getProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["id"]
		resp, err := s.store.GetProjectBySlug(slug)
		if err != nil {
			s.respond(w, r, fmt.Sprintf("Project %s not found", slug), http.StatusNotFound)
			return
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

// TODO: refacto (repetition with updateProject())
func (s *server) createProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(0); err != nil {
			clog.Printlb(err, "FORM-DATA PARSING ERROR")
			s.respond(w, r, "Invalid input", http.StatusBadRequest)
			return
		}

		p, err := NewProjectFromPostForm(r.PostForm)
		if err != nil {
			s.respond(w, r, fmt.Sprint(err), http.StatusBadRequest)
			return
		}

		if err := s.store.InsertProject(p); err != nil {
			s.respond(w, r, fmt.Sprint(err), http.StatusBadRequest)
			return
		}

		s.respond(w, r, p.formatJSON(), http.StatusCreated)
	}
}

// TODO: refacto
func (s *server) updateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.respond(w, r, "Only integers are accepted for this route (Project.ID)", http.StatusBadRequest)
			return
		}

		p0, err := s.store.GetProjectByID(id)
		if err != nil {
			s.respond(w, r, fmt.Sprintf("Project with id %d does not exist", id), http.StatusBadRequest)
			return
		}

		if err := r.ParseMultipartForm(0); err != nil {
			s.respond(w, r, "Invalid input, form-data expected", http.StatusBadRequest)
			return
		}

		p1, err := NewProjectFromPostForm(r.PostForm)
		if err != nil {
			s.respond(w, r, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		p1.ID = p0.ID
		p1.AddedOn = p0.AddedOn
		p1.Update()

		if err = s.store.UpdateProject(p1); err != nil {
			s.respond(w, r, "Project could not be updated", http.StatusBadRequest)
			clog.Printlb(err, "UPDATE ERROR")
			return
		}

		type updateResponse struct {
			Before Project `json:"before"`
			After  Project `json:"after"`
		}

		resp := updateResponse{
			Before: *p0.formatJSON(),
			After:  *p1.formatJSON(),
		}

		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) deleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, err := utl.SplitInt64(mux.Vars(r)["id"], ",")
		if err != nil || len(ids) == 0 {
			s.respond(w, r, "Only comma-separated integers are accepted for this route", http.StatusBadRequest)
			return
		}

		if err = s.store.DeleteManyProjects(ids); err != nil {
			s.respond(w, r, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, "Project deletion successful", http.StatusOK)
	}
}
