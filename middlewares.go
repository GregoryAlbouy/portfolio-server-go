package main

import (
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// authOnly requires for the wrapped HandlerFunc a valid auth token
func (s *server) authOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("APP_JWT_KEY")), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})

		j.HandlerWithNext(w, r, next)
	}
}

// adminOnly wraps a HandlerFunc and blocks it if envitonment variable
// APP_ENV is not set to "admin"
func (s *server) adminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminMode() {
			s.forbidden(w, r)
			return
		}
		next(w, r)
	}
}

// adminOnlyMiddleware blocks any route under a (sub)router that uses it
// (router.Use(adminOnlyMiddleware)) if environment variable APP_ENV is
// not set to "admin"
func (s *server) adminOnlyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdminMode() {
			s.forbidden(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// corsMiddleware sets basic cors to the response header.
// Had to set it manually since I didn't manage to get the gorilla/mux handler
// to work properly on non-GET routes.
func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type")
		h.ServeHTTP(w, r)
	})
}
