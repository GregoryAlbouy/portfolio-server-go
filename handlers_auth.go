package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (s *server) createToken() http.HandlerFunc {
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")
	key := os.Getenv("AUTH_JWT_KEY")

	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		Token string `json:"token"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}

		// Check payload
		if err := s.decode(w, r, &req); err != nil {
			msg := fmt.Sprintf("Cannot parse login body. Err: %v", err)
			s.respond(w, r, responseError{msg}, http.StatusBadRequest)
			return
		}

		// Check credentials
		if req.Username != username || req.Password != password {
			s.respond(w, r, responseError{"Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat":      time.Now().Unix(),
		})
		tokenStr, err := token.SignedString([]byte(key))
		if err != nil {
			msg := fmt.Sprintf("Cannot generate JWT. Err: %v", err)
			s.respond(w, r, responseError{msg}, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, response{tokenStr}, http.StatusOK)
	}
}
