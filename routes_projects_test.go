package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	projectID int64
	projects  []*Project
}

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
	const (
		method = "POST"
		url    = "/api/v1/projects/"
	)

	srv := newServer()
	srv.store = &testStore{}

	p := map[string]string{
		"name":        "Test project",
		"slug":        "test-project",
		"description": "This is a test project.",
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	for k, v := range p {
		err := writer.WriteField(k, v)
		assert.Nil(t, err)
	}

	err := writer.Close()
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/projects/", payload)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	srv.createProject()(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateProjectIntegration(t *testing.T) {
	const (
		method = "POST"
		url    = "/api/v1/projects/"
	)

	srv := newServer()
	srv.store = &testStore{}

	p := map[string]string{
		"name":        "Test project",
		"slug":        "test-project",
		"description": "This is a test project.",
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	for k, v := range p {
		err := writer.WriteField(k, v)
		assert.Nil(t, err)
	}

	err := writer.Close()
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/projects/", payload)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	srv.router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}
