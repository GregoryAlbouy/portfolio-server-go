package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"gregoryalbouy-server-go/db"
	"gregoryalbouy-server-go/utl"
)

const baseURL = "/api/v1"

var routes = utl.RouteMap{
	"/": handleDefaultRoute,
}

func handleDefaultRoute(w http.ResponseWriter, r *http.Request) {
	if string(r.URL.Path) == "/" {
		return
	}

	msg := []byte("Welcome")
	w.Write(msg)
}

func main() {
	var (
		port   = os.Getenv("PORT")
		dburi  = os.Getenv("DB_URI")
		dbname = os.Getenv("DB_NAME")
	)

	db.Connect(dburi, dbname)
	serveRoutes(routes)
	startServer(port)
}

func serveRoutes(routes utl.RouteMap) {
	// for path, handler := range routes {
	// 	http.HandleFunc(path, handler)
	// }

	// for path, handler := range projects.Routes {
	// 	http.HandleFunc(baseURL+path, handler)
	// }
}

func test(w http.ResponseWriter, r *http.Request) {

}

func startServer(port string) {
	const (
		defaultPort = "8080"
		timeout     = 10 * time.Second
	)

	if port == "" {
		port = defaultPort
	}

	fmt.Println("Listening to port " + port)

	router := mux.NewRouter()
	srv := &http.Server{
		Addr:    "127.0.0.1:" + port,
		Handler: router,
	}

	router.HandleFunc("/products", test).
		Host("www.example.com").
		Methods("GET").
		Schemes("http")

	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	log.Fatal("Server cannot listen to port: " + port)
	// }

	srv.ListenAndServe()
}
