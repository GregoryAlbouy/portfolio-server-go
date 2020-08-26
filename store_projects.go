// TODO: Create distinct store for Project methods

package main

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

func (store *dbStore) InsertProject(p *Project) error {
	if store.projectExists(p) {
		return fmt.Errorf("project %s already exists", p.Slug)
	}

	res, err := store.db.NamedExec(queries["project_insert"], p.formatSQL())
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
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

func (store *dbStore) projectExists(p *Project) bool {
	pj, _ := store.GetProjectBySlug(p.Slug)
	return pj != nil
}
