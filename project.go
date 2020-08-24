package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
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
	p.AddedOn = time.Now().Unix()

	return &p
}

// NewProjectFromMap ctakes a map as input and returns a Project struct
// TODO
func NewProjectFromMap(m map[string]interface{}) *Project {
	p := NewProject()
	return p
}

// NewProjectFromForm takes a formdata as input and returns a Project struct
// TODO
func NewProjectFromForm(f url.Values) *Project {
	p := NewProject()
	return p
}

// Init adds current timestamp to AddedOn field if not created
// with NewProject() method
func (p *Project) Init() *Project {
	p.AddedOn = timestamp()
	return p
}

// Update updates the value of EditedOn field with the current time
func (p *Project) Update() *Project {
	p.EditedOn = timestamp()
	return p
}

// SetTagstr sets TagStr field from Tags value
func (p *Project) setTagstr() *Project {
	p.Tagstr = strings.Join(p.Tags, ",")
	return p
}

// SetTags sets Tags fields from TagStr value (string with comma-separated values)
func (p *Project) setTags() *Project {
	p.Tags = strings.Split(p.Tagstr, ",")
	return p
}

// IsValid checks whether a project contains the required fields
func (p *Project) IsValid() bool {
	return p.Name != "" && p.Slug != "" && p.Description != ""
}

func (p *Project) formatJSON() *Project {
	p.setTags()
	return p
}

func (p *Project) formatSQL() *Project {
	p.setTagstr()
	return p
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
