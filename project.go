package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ProjectList struct
type ProjectList []*Project

func (pl ProjectList) toJSON() *ProjectList {
	for _, p := range pl {
		p.toJSON()
		fmt.Println(p)
	}
	return &pl
}

// Project struct
type Project struct {
	ID          int64   `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Slug        string  `db:"slug" json:"slug"`
	Description string  `db:"description" json:"description"`
	Tags        Taglist `db:"-" json:"tags"`
	TagStr      string  `db:"tags" json:"-"`
	Image       string  `db:"image" json:"image,omitempty"`
	Repo        string  `db:"repo" json:"repo,omitempty"`
	Demo        string  `db:"demo" json:"demo,omitempty"`
	AddedOn     int64   `db:"added_on" json:"added_on"`
	EditedOn    int64   `db:"edited_on" json:"edited_on,omitempty"`
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

// SetTagStr sets TagStr field from Tags value
func (p *Project) setTagStr() {
	p.TagStr = strings.Join(p.Tags, ",")
}

// SetTags sets Tags fields from TagStr value (string with comma-separated values)
func (p *Project) setTags() {
	p.Tags = strings.Split(p.TagStr, ",")
}

// IsValidForInsertion checks whether a project contains the required fields
func (p *Project) IsValidForInsertion() bool {
	return p.Name != "" && p.Slug != "" && p.Description != ""
}

// func (p *Project) isValidForCreate() bool {
// 	return p.isValidForUpdate() && !p.alreadyInDB()
// }

// func (p *Project) alreadyInDB() bool {
// 	return false
// }

func (p *Project) toJSON() {
	p.setTags()
}

func (p *Project) toDB() {
	p.setTagStr()
}

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
