package main

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")

	pj := s.router.PathPrefix("/api/v1/projects/").Subrouter()
	pj.HandleFunc("/test/", s.handleProjectList()).Methods("GET")
	pj.HandleFunc("/{id:[0-9]+}", s.handleProjectDetail()).Methods("GET")
	pj.HandleFunc("/", s.handleProjectCreation()).Methods("POST")
	pj.HandleFunc("/{id:[0-9]+}", s.handleProjectUpdate()).Methods("PATCH")
	pj.HandleFunc("/{id}", s.handleProjectDeletion()).Methods("DELETE")
}
