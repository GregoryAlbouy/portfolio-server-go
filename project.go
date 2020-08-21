package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Project struct
type Project struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Slug        string  `db:"slug"`
	Description string  `db:"description"`
	Tags        Taglist `db:"-"`
	TagStr      string  `db:"tags"`
	Image       string  `db:"image"`
	Repo        string  `db:"repo"`
	Demo        string  `db:"demo"`
	AddedOn     int64   `db:"added_on"`
	EditedOn    int64   `db:"edited_on"`
}

// Taglist type is a slice of strings that prints by joining them with a comma
type Taglist []string

func (tl Taglist) String() string {
	var str string
	n := len(tl) - 1

	for i, v := range tl {
		sep := ","
		if i == n {
			sep = ""
		}

		str += v + sep
	}

	return str
}

// NewTaglist create a Taglist from a string with comma-separated values.
// Can be used for practical retrieve from database.
func NewTaglist(str string) (tl Taglist) {
	tl = strings.Split(str, ",")
	return
}

// NewProject initializes Project struct with current time as AddedOn value
func NewProject() *Project {
	p := &Project{}
	p.AddedOn = time.Now().Unix()

	return p
}

// Init adds current timestamp to AddedOn field if not created
// with NewProject() method
func (p *Project) Init() {
	p.AddedOn = timestamp()
}

// Update updates the value of EditedOn field with the current time
func (p *Project) Update() {
	p.EditedOn = timestamp()
}

// IsValidForInsertion checks whether a project contains the required fields
func (p *Project) IsValidForInsertion() bool {
	return p.Name != "" && p.Slug != "" && p.Description != ""
}

// func (p *Project) isValidForCreate() bool {
// 	return p.isValidForUpdate() && !p.alreadyInDB()
// }

// func (p *Project) alreadyInDB() bool {
// 	return len()
// }

func (p Project) String() string {
	j, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return string(j) + "\n"
}

func timestamp() int64 {
	return time.Now().Unix()
}
