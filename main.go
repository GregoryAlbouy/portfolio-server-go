package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// func tests(srv *server) {
// 	p1 := Project{
// 		Name:        "Codemon",
// 		Slug:        "codemon",
// 		Description: "School project",
// 		Repo:        "https://github.com/gregoryalbouy/codemon",
// 		Tags: Taglist{
// 			"OO JS",
// 			"ES6",
// 			"JS Modules",
// 			"3D CSS Animation",
// 		},
// 		// Tags: NewTaglist("Node.js,Golang,Python"),
// 	}
// 	p1.Init()

// 	createProject := func(p *Project) {
// 		if err := srv.store.CreateProject(&p1); err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	getProject := func(p *Project) {
// 		project, err := srv.store.GetProjectBySlug(p.Slug)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println(project)
// 	}

// 	getAllProjects := func() {
// 		pjs, err := srv.store.GetProjects()
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println(pjs)
// 	}

// 	clearDB := func() {
// 		if err := srv.store.Clear(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}

// 	dropDB := func() {
// 		if err := srv.store.Drop(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}

// 	f1 := []func(){dropDB, clearDB, getAllProjects}
// 	f2 := []func(*Project){getProject, createProject}
// 	fmt.Sprintf("%v%v", f1, f2)

// 	// createProject(&p1)
// 	getAllProjects()
// }

func main() {
	fmt.Println("Starting.")

	if err := run(); err != nil {
		// fmt.Fprintf(os.Stderr, "%s\n", err)
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() (err error) {
	srv := newServer()
	srv.store = &dbStore{}

	err = srv.store.Open()
	if err != nil {
		return
	}
	defer srv.store.Close()

	http.HandleFunc("/", srv.serveHTTP)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

	// tests(srv)
	return
}
