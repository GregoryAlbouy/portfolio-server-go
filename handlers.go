package main

import (
	"fmt"
	"gregoryalbouy-server-go/clog"
	"gregoryalbouy-server-go/utl"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}
}

func (s *server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["id"]
		resp, err := s.store.GetUserByUsername(name)
		if err != nil {
			s.respond(w, r, fmt.Sprintf("User %s not found", name), http.StatusNotFound)
			return
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.respond(w, r, "Only integers are accepted for this route", http.StatusBadRequest)
			return
		}

		if err = s.store.DeleteUser(id); err != nil {
			s.respond(w, r, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, "User deletion successful", http.StatusOK)
	}
}

func (s *server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &User{}

		if err := s.decodeRequest(r, u); err != nil {
			s.respond(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		if !u.IsValid() {
			s.respond(w, r, "Username must be >= 3 char and password >= 8 char", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.RawPassword), 10)
		if err != nil {
			clog.Printlb(err, "HASH ERROR - createUser()")
			s.respond(w, r, "Internal error", http.StatusInternalServerError)
		}
		u.Password = string(hash)

		if err := s.store.InsertUser(u); err != nil {
			clog.Printlb(err, "USER INSERTION ERROR")
			s.respond(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		s.respond(w, r, u.Safe(), http.StatusCreated)
	}
}

func (s *server) getUserList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.GetUserList()
		if err != nil {
			s.respond(w, r, "No result", http.StatusNoContent)
			return
		}

		s.respond(w, r, users, http.StatusOK)
	}
}

func (s *server) createToken() http.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := &User{}

		// Check payload
		if err := s.decodeRequest(r, u); err != nil {
			s.respond(w, r, responseError{"Invalid payload"}, http.StatusBadRequest)
			return
		}

		// Check credentials
		if err := s.AuthenticateUser(u); err != nil {
			s.respond(w, r, responseError{"Invalid credentials"}, http.StatusBadRequest)
			return
		}

		// Generate token
		token, err := tokenFromUser(u)
		if err != nil {
			s.respond(w, r, responseError{"Cannot generate token"}, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, response{token}, http.StatusOK)
	}
}

func (s *server) AuthenticateUser(input *User) error {
	// Check username
	dbu, err := s.store.GetUserByUsername(input.Username)
	if err != nil {
		clog.Printlb(err, "AUTHENTICATION ERROR - Username")
		return err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(dbu.Password), []byte(input.RawPassword))
	if err != nil {
		clog.Printlb(err, "AUTHENTICATION ERROR - Password")
		return err
	}

	return nil
}

func tokenFromUser(u *User) (string, error) {
	key := []byte(os.Getenv("APP_JWT_KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":      time.Now().Unix(),
	})
	tstr, err := t.SignedString(key)
	if err != nil {
		clog.Printlb(err, "TOKEN GENERATION ERROR")
		return "", err
	}

	return tstr, nil
}

func (s *server) postMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := NewMessage()
		if err := s.decodeRequest(r, m); err != nil {
			clog.Printlb(err, clog.Red("DECODE ERROR"))
			s.respond(w, r, err.Error(), http.StatusBadRequest)
			return
		}
		m.IP = utl.RequestIP(r)

		if !m.Valid() {
			clog.Printlb(m, clog.Red("MESSAGE VALIDATION ERROR"))
			s.respond(w, r, "Invalid message", http.StatusBadRequest)
			return
		}

		if err := s.store.InsertMessage(m); err != nil {
			clog.Printlb(err, clog.Red("INSERT ERROR"))
			s.respond(w, r, "Internal error", http.StatusInternalServerError)
			return
		}

		clog.Printlb(m, clog.Green("MESSAGE POSTED"))
		s.respond(w, r, nil, http.StatusCreated)
	}
}

func (s *server) getMessageList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ml, err := s.store.GetMessageList()
		if err != nil {
			s.respond(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, ml, http.StatusOK)
	}
}

func (s *server) deleteMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

		if err != nil {
			// If id can't be parsed to int, treat it as an email
			n, err := s.store.DeleteMessagesByEmail(idstr)
			if err != nil {
				s.respond(w, r, "Message could not be deleted", http.StatusInternalServerError)
			}
			s.respond(w, r, fmt.Sprintf("%d messages deleted", n), http.StatusOK)
			return
		}

		if err := s.store.DeleteMessageByID(id); err != nil {
			fmt.Println(err)
			s.respond(w, r, fmt.Sprintf("Message with ID %d does not exist", id), http.StatusBadRequest)
			return
		}

		s.respond(w, r, nil, http.StatusOK)
	}
}
