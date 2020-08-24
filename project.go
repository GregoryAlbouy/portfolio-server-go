package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

// ProjectList struct
type ProjectList []*Project

func (pl *ProjectList) formatJSON() *ProjectList {
	for _, p := range *pl {
		p.formatJSON()
	}
	return pl
}

// Project struct
type Project struct {
	ID          int64    `db:"id" json:"id"`
	Name        string   `db:"name" json:"name"`
	Slug        string   `db:"slug" json:"slug"`
	Description string   `db:"description" json:"description"`
	Tags        []string `db:"-" json:"tags"`
	Tagstr      string   `db:"tags" json:"-"`
	Image       string   `db:"image" json:"image,omitempty"`
	Repo        string   `db:"repo" json:"repo,omitempty"`
	Demo        string   `db:"demo" json:"demo,omitempty"`
	IsHidden    bool     `db:"is_hidden" json:"is_hidden" schema:"is_hidden"`
	AddedOn     int64    `db:"added_on" json:"added_on"`
	EditedOn    int64    `db:"edited_on" json:"edited_on,omitempty"`
}

// NewProject initializes Project struct with current time as AddedOn value
func NewProject() *Project {
	p := Project{}
	p.setAddedOn()

	return &p
}

// NewProjectFromMap takes a map as input and returns a Project struct
// TODO
func NewProjectFromMap(m map[string]interface{}) *Project {
	return NewProject()
}

// NewProjectFromPostForm takes a formdata as input and returns a Project struct
// TODO
func NewProjectFromPostForm(pf url.Values) (*Project, error) {
	p := NewProject()
	decoder := schema.NewDecoder()

	if err := decoder.Decode(p, pf); err != nil {
		fmt.Println(err)
		return nil, errors.New("Invalid project input. Form-data expected")
	}

	if !p.IsValid() {
		return nil, errors.New("Incomplete project data")
	}

	return p, nil
}

// // Assign assigns the input *Project fields values to the current *Project
// // except for ID, AddedOn, EditedOn. EditedOn receives a timestamp
// func (p *Project) Assign(p1 *Project) {
// 	if p.Tagstr == "" {
// 		p.setTagstr()
// 	}
// 	if p1.Tagstr == "" {
// 		p1.setTagstr()
// 	}
// 	p.Name = p1.Name
// 	p.Slug = p1.Slug
// 	p.Description = p1.Description
// 	p.Tagstr = p1.Tagstr
// 	p.Image = p1.Image
// 	p.Repo = p1.Repo
// 	p.Demo = p1.Demo
// 	p.IsHidden = p1.IsHidden
// 	p.Update()
// }

// Init adds current timestamp to AddedOn field if not created
// with NewProject() method
func (p *Project) Init() *Project {
	return p.setAddedOn()
}

// Update updates the value of EditedOn field with the current time
func (p *Project) Update() *Project {
	return p.setEditedOn()
}

// SetTagstr sets TagStr field from Tags value
func (p *Project) setTagstr() *Project {
	p.Tagstr = strings.Join(p.Tags, ",")
	return p
}

// SetTags sets Tags fields from TagStr value (string with comma-separated values)
func (p *Project) setTags() *Project {
	if p.Tagstr != "" {
		p.Tags = strings.Split(p.Tagstr, ",")
	} else {
		p.Tags = []string{}
	}
	return p
}

// IsValid checks whether a project contains the required fields
func (p *Project) IsValid() bool {
	return p.Name != "" && p.Slug != "" && p.Description != "" && p.AddedOn != 0
}

func (p *Project) formatJSON() *Project {
	return p.setTags()
}

func (p *Project) formatSQL() *Project {
	if p.AddedOn == 0 {
		p.setAddedOn()
	}
	return p.setTagstr()
}

func (p Project) String() string {
	j, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return string(j) + "\n"
}

func (p *Project) setAddedOn() *Project {
	p.AddedOn = timestamp()
	return p
}

func (p *Project) setEditedOn() *Project {
	p.EditedOn = timestamp()
	return p
}

func timestamp() int64 {
	return time.Now().Unix()
}
