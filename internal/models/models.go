package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Admin    Bool      `json:"admin"`
}

func (u *User) Copy() *User {
	user := *u
	return &user
}

func (u *User) Patch(new User) {
	if new.Username != "" {
		u.Username = new.Username
	}
	if new.Password != "" {
		u.Password = new.Password
	}
	if new.Email != "" {
		u.Email = new.Email
	}
	if new.Admin.Valid() {
		u.Admin = new.Admin
	}
}
