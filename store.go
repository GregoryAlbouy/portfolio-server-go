package main

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var projectSchema = `
CREATE TABLE IF NOT EXISTS project
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT NOT NULL,
	tags TEXT,
	image TEXT,
	repo TEXT,
	demo TEXT,
	is_hidden INTEGER,
	added_on INTEGER,
	edited_on INTEGER
)
`

// Store interface
type Store interface {
	Open() error
	Close() error

	GetProjectList() (ProjectList, error)
	GetProjectBySlug(string) (*Project, error)
	UpdateProjectBySlug(string) (*Project, error)
	CreateProject(*Project) error
	Clear() error
	Drop() error
}

type dbStore struct {
	db *sqlx.DB
}

func (store *dbStore) Open() error {
	db, err := sqlx.Connect("sqlite3", ".db")
	if err != nil {
		return err
	}

	fmt.Println("Connected to DB")

	db.MustExec(projectSchema)
	store.db = db

	return nil
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) GetProjectList() (pl ProjectList, err error) {
	err = store.db.Select(&pl, "SELECT * FROM project ORDER BY added_on DESC")
	return
}

func (store *dbStore) GetProjectBySlug(slug string) (*Project, error) {
	var res []*Project
	err := store.db.Select(&res, fmt.Sprintf("SELECT * FROM project WHERE slug='%s' LIMIT 1", slug))
	if err != nil {
		return nil, errors.New("Store error: SELECT")
	}
	if len(res) == 0 {
		return nil, errors.New("not found")
	}
	return res[0], nil
}

func (store *dbStore) GetProjectByID(id int64) (*Project, error) {
	var res []*Project
	err := store.db.Select(&res, fmt.Sprintf("SELECT * FROM project WHERE id='%d' LIMIT 1", id))
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

func (store *dbStore) UpdateProjectBySlug(slug string) (*Project, error) {
	rowx := store.db.QueryRowx("SELECT * FROM project")
	res := map[string]interface{}{}
	err := rowx.MapScan(res)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
	return &Project{}, nil
}

func (store *dbStore) CreateProject(p *Project) error {
	if store.projectExists(p.Slug) {
		fmt.Printf("Project %s already exists\n", p.Slug)
		return nil
	}

	p.formatSQL()

	// res, err := store.db.Exec("INSERT INTO project (name, slug, description, tags, image, repo, demo, added_on, edited_on) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", p.Name, p.Slug, p.Description, p.Tagstr, p.Image, p.Repo, p.Demo, p.AddedOn, p.EditedOn)
	// if err != nil {
	// 	return err
	// }

	res, err := store.db.NamedExec("INSERT INTO project (name, slug, description, tags, image, repo, demo, added_on, edited_on) VALUES (:name, :slug, :description, :tags, :image, :repo, :demo, :added_on, :edited_on)", p)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("Project %s added to database\n", p.Slug)
	return nil
}

func (store *dbStore) Clear() error {
	_, err := store.db.Exec("DELETE FROM project")
	if err != nil {
		return err
	}
	fmt.Println("Table project cleared")

	return nil
}

func (store *dbStore) Drop() error {
	_, err := store.db.Exec("DROP TABLE IF EXISTS project")
	if err != nil {
		return err
	}
	fmt.Println("Table project dropped")
	return nil
}

func (store *dbStore) projectExists(slug string) bool {
	p, err := store.GetProjectBySlug(slug)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(slug)
	return p != nil
}
