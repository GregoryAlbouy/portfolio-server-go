package main

import "fmt"

type testStore struct {
	pcount   int64
	projects []*Project
	ucount   int64
	users    []*User
}

func (t testStore) InsertProject(p *Project) error {
	t.pcount++
	p.ID = t.pcount
	t.projects = append(t.projects, p)
	return nil
}

func (t *testStore) GetProjectList() (ProjectList, error) {
	return t.projects, nil
}

func (t testStore) GetProjectByID(id int64) (*Project, error) {
	for _, p := range t.projects {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Project ID %d not found", id)
}

func (t testStore) GetProjectBySlug(slug string) (*Project, error) {
	for _, p := range t.projects {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Project %s not found", slug)
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
func (t testStore) Clear(table string) error {
	return nil
}
func (t testStore) Drop(table string) error {
	return nil
}
func (t testStore) InsertUser(u *User) error {
	return nil
}
func (t testStore) DeleteUser(id int64) error {
	return nil
}
func (t testStore) GetUserByUsername(username string) (*User, error) {
	return nil, nil
}
func (t testStore) GetUserList() ([]*User, error) {
	return nil, nil
}
func (t testStore) UserExists(*User) bool {
	return true
}
func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}
