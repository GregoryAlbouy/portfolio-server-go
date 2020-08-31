package main

import (
	"encoding/json"
	"fmt"
	"gregoryalbouy-server-go/clog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	addr    string
	port    string
	router  *mux.Router
	store   Store
	started time.Time
}

func newServer() *server {
	r := mux.NewRouter().StrictSlash(true)

	// r.Use(handlers.CORS(
	// 	handlers.AllowedHeaders([]string{"Content-Type"}),
	// 	handlers.AllowedOrigins([]string{"*"}),
	// 	handlers.AllowCredentials(),
	// ))

	s := &server{
		started: time.Now(),
		router:  r,
	}

	s.router.Use(corsMiddleware)

	s.routes()
	return s
}

func (s *server) setPortWithDefault(defaultPort string) *server {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	s.port = port
	s.addr = ":" + port

	fmt.Println("Port", clog.Blue(s.port))
	return s
}

func (s *server) attachStore(store *dbStore) *server {
	s.store = store
	return s
}

func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("JSON Marshal error: %v\n", err)
	}

	w.Write(resp)
}

func (s *server) forbidden(w http.ResponseWriter, r *http.Request) {
	s.respond(w, r, "Access forbidden", http.StatusForbidden)
}

func (s *server) decodeRequest(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("Cannot parse request body: %s", err)
	}
	return nil
}

func (s *server) isAdminMode() bool {
	env := os.Getenv("APP_ENV")
	return env == "admin" || env == "dev"
}

func (s *server) serve() error {
	errChan := make(chan error)
	go func() {
		if err := http.ListenAndServe(s.addr, s.router); err != nil {
			errChan <- err
		}
	}()
	s.printStatus()
	err := <-errChan
	return err
}

func (s *server) printStatus() {
	fmt.Printf("%s (%s) http://127.0.0.1%s \n", clog.Green("Ready"), time.Since(s.started), s.addr)
}
