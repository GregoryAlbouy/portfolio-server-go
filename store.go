package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"gregoryalbouy-server-go/clog"
)

var queries = map[string]string{
	"table_project_create": "",
	"table_project_clear":  "",
	"table_project_drop":   "",
	"project_insert":       "",
	"project_all":          "",
	"project_ids":          "",
	"project_by_id":        "",
	"project_by_slug":      "",
	"project_update":       "",
	"project_delete":       "",
	"project_delete_many":  "",
}

// Store interface
type Store interface {
	Open() error
	Close() error

	GetProjectList() (ProjectList, error)
	GetProjectByID(int64) (*Project, error)
	GetProjectBySlug(string) (*Project, error)
	UpdateProject(*Project) error
	CreateProject(*Project) error
	DeleteProject(int64) error
	DeleteManyProjects([]int64) error
	Clear() error
	Drop() error
}

type dbStore struct {
	db *sqlx.DB
}

func init() {
	// Read raw SQL queries from files and store them into query map
	for k := range queries {
		query, err := ioutil.ReadFile(fmt.Sprintf("queries/%s.sql", k))
		if err != nil {
			clog.Fatallb(err, "STORE (init)")
		}
		queries[k] = string(query)
	}
}

func (store *dbStore) Open() error {
	t0 := time.Now()
	db, err := sqlx.Connect("sqlite3", ".db")
	if err != nil {
		return err
	}

	db.MustExec(queries["table_project_create"])
	store.db = db

	fmt.Printf("DB connection %s (%s)\n", clog.Green("OK"), time.Since(t0))
	return nil
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) GetProjectList() (pl ProjectList, err error) {
	if err := store.db.Select(&pl, queries["project_all"]); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}

func (store *dbStore) GetProjectByID(id int64) (*Project, error) {
	p := Project{}
	if err := store.db.QueryRowx(queries["project_by_id"], id).StructScan(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (store *dbStore) GetProjectBySlug(slug string) (*Project, error) {
	p := Project{}
	if err := store.db.QueryRowx(queries["project_by_slug"], slug).StructScan(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (store *dbStore) UpdateProject(p *Project) error {
	_, err := store.db.NamedExec(queries["project_update"], p.formatSQL())
	if err != nil {
		fmt.Println(queries["project_update"])
		return err
	}

	return nil
}

func (store *dbStore) CreateProject(p *Project) error {
	if store.projectExists(p) {
		return fmt.Errorf("project %s already exists", p.Slug)
	}

	_, err := store.db.NamedExec(queries["project_insert"], p.formatSQL())
	if err != nil {
		return err
	}

	return nil
}

func (store *dbStore) DeleteProject(id int64) error {
	_, err := store.db.Exec(queries["project_delete"], id)
	if err != nil {
		return errors.New("Project does not exist")
	}
	return nil
}

func (store *dbStore) DeleteManyProjects(ids []int64) error {
	ok, err := store.projectIdsExist(ids)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Project does not exist")
	}

	q, args, err := sqlx.In(queries["project_delete_many"], ids)
	if err != nil {
		return errors.New("Internal error")
	}

	q = store.db.Rebind(q)
	_, err = store.db.Exec(q, args...)
	if err != nil {
		return errors.New("Internal error")
	}

	return nil
}

func (store *dbStore) Clear() error {
	_, err := store.db.Exec(queries["table_project_clear"])
	if err != nil {
		return err
	}
	fmt.Println("Table project cleared")

	return nil
}

func (store *dbStore) Drop() error {
	_, err := store.db.Exec(queries["table_project_drop"])
	if err != nil {
		return err
	}
	fmt.Println("Table project dropped")
	return nil
}

func (store *dbStore) projectExists(p *Project) bool {
	pj, _ := store.GetProjectBySlug(p.Slug)
	return pj != nil
}

func (store *dbStore) projectIdsExist(ids []int64) (bool, error) {
	pIds := []int64{}

	q, args, err := sqlx.In(queries["project_ids"], ids)
	if err != nil {
		return false, errors.New("Internal error")
	}

	q = store.db.Rebind(q)
	err = store.db.Select(&pIds, q, args...)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("Internal error")
	}

	return len(ids) == len(pIds), nil
}
