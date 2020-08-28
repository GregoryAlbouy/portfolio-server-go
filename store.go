package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"gregoryalbouy-server-go/clog"
)

var queries = map[string]string{
	"table_project_create": "",
	"table_user_create":    "",
	"table_message_create": "",
	"project_insert":       "",
	"project_all":          "",
	"project_ids":          "",
	"project_by_id":        "",
	"project_by_slug":      "",
	"project_update":       "",
	"project_delete":       "",
	"project_delete_many":  "",
	"user_all":             "",
	"user_by_username":     "",
	"user_delete":          "",
	"user_insert":          "",
	"message_insert":       "",
}

// Store interface
type Store interface {
	Open() error
	Close() error

	Clear(string) error
	Drop(string) error

	GetProjectList() (ProjectList, error)
	GetProjectByID(int64) (*Project, error)
	GetProjectBySlug(string) (*Project, error)
	UpdateProject(*Project) error
	InsertProject(*Project) error
	DeleteProject(int64) error
	DeleteManyProjects([]int64) error

	GetUserList() ([]*User, error)
	InsertUser(*User) error
	GetUserByUsername(string) (*User, error)
	DeleteUser(int64) error
	UserExists(*User) bool

	GetMessageList() ([]*Message, error)
	InsertMessage(*Message) error
	DeleteMessageByID(int64) error
	DeleteMessagesByEmail(string) (int64, error)
}

type dbStore struct {
	db *sqlx.DB
}

func init() {
	// Read raw SQL queries from files and store them into query map
	for k := range queries {
		query, err := ioutil.ReadFile(fmt.Sprintf("queries/%s.sql", k))
		if err != nil {
			clog.Fatallb(err, "STORE init()")
		}
		queries[k] = string(query)
	}
}

func (store *dbStore) Open() error {
	path := os.Getenv("APP_DB_PATH")
	status := clog.Green("OK")

	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		status = clog.Red("Error")
	}
	store.db = db
	store.createTables()

	fmt.Printf("DB path %s status %s\n", clog.Blue("/"+path), status)
	return err
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) createTables() {
	store.db.MustExec(queries["table_project_create"])
	store.db.MustExec(queries["table_user_create"])
	store.db.MustExec(queries["table_message_create"])
}

func (store *dbStore) Clear(table string) error {
	query := fmt.Sprintf("DELETE FROM %s", table)
	_, err := store.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Printf("Table %s cleared\n", table)
	return nil
}

func (store *dbStore) Drop(table string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	_, err := store.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Printf("Table %s dropped\n", table)
	return nil
}

func (store *dbStore) GetUserList() (ul []*User, err error) {
	if err := store.db.Select(&ul, queries["user_all"]); err != nil {
		return nil, err
	}
	return
}

func (store *dbStore) GetUserByUsername(username string) (*User, error) {
	u := User{}
	if err := store.db.QueryRowx(queries["user_by_username"], username).StructScan(&u); err != nil {
		return nil, err
	}
	if u.Username == "" {
		return nil, fmt.Errorf("User %s does not exist", username)
	}
	return &u, nil
}

func (store *dbStore) InsertUser(u *User) error {
	if store.UserExists(u) {
		return fmt.Errorf("user %s already exists", u.Username)
	}

	res, err := store.db.NamedExec(queries["user_insert"], u)
	if err != nil {
		return err
	}

	u.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// TODO: check if user exists first
func (store *dbStore) DeleteUser(id int64) error {
	_, err := store.db.Exec(queries["user_delete"], id)
	if err != nil {
		return fmt.Errorf("User with id %d does not exist", id)
	}
	return nil
}

func (store *dbStore) UserExists(u *User) bool {
	user, _ := store.GetUserByUsername(u.Username)
	return user != nil
}

func (store *dbStore) GetMessageList() (ml []*Message, err error) {
	err = store.db.Select(&ml, "SELECT * FROM message ORDER BY date DESC")
	return
}

func (store *dbStore) InsertMessage(m *Message) error {
	_, err := store.db.NamedExec(queries["message_insert"], m)
	return err
}

func (store *dbStore) DeleteMessageByID(id int64) error {
	_, err := store.db.Exec("DELETE FROM message WHERE id=?", id)
	return err
}

func (store *dbStore) DeleteMessagesByEmail(email string) (int64, error) {
	res, err := store.db.Exec("DELETE FROM message WHERE email=?", email)
	if err != nil {
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}
