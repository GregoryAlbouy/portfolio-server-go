package main

const (
	minUsernameLength = 3
	minPasswordLength = 8
)

// User represents... a user...
type User struct {
	ID          int64  `db:"id" json:"id"`
	Username    string `db:"username" json:"username"`
	Password    string `db:"password" json:"-"`
	RawPassword string `db:"-" json:"password,omitempty"`
}

// type UserInput struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// IsValid makes sure Username and RawPassword are long enough before database insertion
func (u *User) IsValid() bool {
	return len(u.Username) >= minUsernameLength &&
		len(u.RawPassword) >= minPasswordLength
}

// Safe removes sensitive data before server response
func (u *User) Safe() *User {
	u.Password = "" // Should be useless with "-" struct tag but anyways
	u.RawPassword = ""
	return u
}
