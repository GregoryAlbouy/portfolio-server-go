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

type testStore struct {
	projectID int64
	projects  []*Project
}

type testRequest struct {
	method  string
	target  string
	headers map[string]string
	body    io.Reader
}

var ct = map[string]string{
	"json": "application/json",
}

var path = map[string]string{
	"auth":    "/auth/",
	"project": "/api/v1/projects/",
}

var srv = newTestServer()
var testProject = validProject()
var w = httptest.NewRecorder()

func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}

func (t *testStore) GetProjectList() (ProjectList, error) {
	return t.projects, nil
}

func (t testStore) GetProjectBySlug(slug string) (*Project, error) {
	for _, p := range t.projects {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, nil
}

func (t testStore) CreateProject(p *Project) error {
	t.projectID++
	p.ID = t.projectID
	t.projects = append(t.projects, p)
	return nil
}
func (t testStore) GetProjectByID(id int64) (*Project, error) {
	return nil, nil
}
func (t testStore) UpdateProject(p *Project) error {
	return nil
}
func (t testStore) DeleteManyProjects(ids []int64) error {
	return nil
}
func (t testStore) DeleteProject(id int64) error {
	return nil
}
func (t testStore) Clear() error {
	return nil
}
func (t testStore) Drop() error {
	return nil
}

func TestCreateProjectUnit(t *testing.T) {
	fd, ctype, err := createFormData(testProject)
	assert.Nil(t, err)

	r := (&testRequest{
		method: "POST",
		target: path["project"],
		headers: map[string]string{
			"Content-Type": ctype,
		},
		body: fd,
	}).Set()

	srv.createProject()(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func testAuthRoute(t *testing.T) {

}

func TestCreateProjectIntegration(t *testing.T) {
	fd, ctype, err := createFormData(testProject)
	assert.Nil(t, err)

	token, err := validToken()
	assert.Nil(t, err)

	r := (&testRequest{
		method: "POST",
		target: path["project"],
		headers: map[string]string{
			"Authorization": authHeader(token),
			"Content-Type":  ctype,
		},
		body: fd,
	}).Set()

	srv.router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func credentialsPayload() (io.Reader, error) {
	b, err := json.Marshal(User{
		Username: os.Getenv("AUTH_USERNAME"),
		Password: os.Getenv("AUTH_PASSWORD"),
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func validToken() (string, error) {
	const (
		URL   = "http://localhost:8080/auth/"
		ctype = "application/json"
	)

	type responseToken struct {
		Token string `json:"token"`
	}
	var rt responseToken

	payload, err := credentialsPayload()
	if err != nil {
		return "", err
	}

	resp, err := http.Post(URL, ctype, payload)
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(resp.Body).Decode(&rt); err != nil {
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

func (tr *testRequest) Set() *http.Request {
	r := httptest.NewRequest(tr.method, tr.target, tr.body)

	for k, v := range tr.headers {
		r.Header.Set(k, v)
	}

	return r
}
