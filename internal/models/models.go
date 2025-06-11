package models

import (
	"errors"
	"fmt"

	"encoding/json"
)

type Bool struct {
	val bool
	ok  bool
}

func (b *Bool) Value() bool { return b.val }
func (b *Bool) Valid() bool { return b.ok }

func (b *Bool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	(*b).ok = asString == ""
	if asString == "true" {
		(*b).val = true
	} else if asString == "false" {
		(*b).val = false
	} else {
		return errors.New(
			fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString),
		)
	}

	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.val {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}

type User struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Admin    Bool   `json:"admin,omitzero"`
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
		u.Password = new.Password
	}
	if new.Admin.Valid() {
		u.Admin = new.Admin
	}
}

func (u User) MarshalJSON() ([]byte, error) {
	user := struct {
		Id       int64  `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Admin    Bool   `json:"admin,omitzero"`
	}{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Admin:    u.Admin,
	}

	return json.Marshal(user)
}
