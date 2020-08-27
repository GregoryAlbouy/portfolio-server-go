// TODO: Massive refacto
// This is a huge mess

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testRequest describes an http request
type testRequest struct {
	method  string
	target  string
	headers map[string]string
	body    io.Reader
}

// set method prepares a *testRequest and returns a test writer and a *http.Request
func (tr *testRequest) set() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(tr.method, tr.target, tr.body)

	for k, v := range tr.headers {
		r.Header.Set(k, v)
	}

	return w, r
}

var path = map[string]string{
	"token":    "/token/",
	"projects": "/api/v1/projects/",
}

var srv = newTestServer()
var testProject = validProject()

func TestCreateProjectUnit(t *testing.T) {
	fd, ctype, err := createFormData(testProject)
	assert.Nil(t, err)

	w, r := (&testRequest{
		method: "POST",
		target: path["projects"],
		headers: map[string]string{
			"Content-Type": ctype,
		},
		body: fd,
	}).set()

	srv.createProject()(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func testAuthRoute(t *testing.T) {

}

func TestCreateProjectIntegration(t *testing.T) {
	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	fmt.Println(validToken())

	fd, ctype, err := createFormData(testProject)
	assert.Nil(t, err)

	token, err := validToken()
	assert.Nil(t, err)

	w, r := (&testRequest{
		method: "POST",
		target: path["projects"],
		headers: map[string]string{
			"Authorization": authHeader(token),
			"Content-Type":  ctype,
		},
		body: fd,
	}).set()

	srv.router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func credentialsPayload() (io.Reader, error) {
	b, err := json.Marshal(User{
		Username: os.Getenv("TEST_AUTH_VALID_USERNAME"),
		Password: os.Getenv("TEST_AUTH_VALID_PASSWORD"),
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func validToken() (string, error) {
	payload, err := credentialsPayload()
	if err != nil {
		return "", err
	}

	type responseToken struct {
		Token string `json:"token"`
	}
	var rt responseToken

	w, r := (&testRequest{
		method: "POST",
		target: path["token"],
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		body: payload,
	}).set()

	srv.createToken()(w, r)

	if err := json.NewDecoder(w.Body).Decode(&rt); err != nil {
		return "", err
	}

	return rt.Token, nil
}

func authHeader(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func createFormData(m map[string]string) (*bytes.Buffer, string, error) {
	fd := &bytes.Buffer{}
	w := multipart.NewWriter(fd)

	for k, v := range m {
		if err := w.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return fd, w.FormDataContentType(), nil
}

func newTestServer() *server {
	srv := newServer()
	srv.store = &testStore{}
	return srv
}

func validProject() map[string]string {
	return map[string]string{
		"name":        "Test project",
		"slug":        "test-project",
		"description": "This is a test project.",
	}
}
