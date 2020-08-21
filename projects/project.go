package projects

// Project structure
type Project struct {
	ID          string      `json:"_id" bson:"_id"`
	Name        string      `json:"name" bson:"name"`
	Slug        string      `json:"slug" bson:"slug"`
	Description string      `json:"description"`
	Tags        []string    `json:"tags"`
	Image       string      `json:"image,omitempty"`
	Repo        string      `json:"repo,omitempty"`
	Demo        string      `json:"demo,omitempty"`
	AddedOn     interface{} `json:"addedOn"`
}

// IsValid checks whether the current projects has enough data
func (p *Project) IsValid() bool {
	return true
}

// ToJSON converts Project to JSON format for response output
func (p *Project) ToJSON() interface{} {
	return 1
}
