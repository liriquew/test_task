package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func GetDefaultAdmin() *User {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &User{
		Id:       id,
		Username: "admin",
		Email:    "admin@admin.ru",
		Password: "admin",
		Admin: Bool{
			Val: true,
			Ok:  true,
		},
	}
}

type Bool struct {
	Val bool
	Ok  bool
}

func (b *Bool) Value() bool { return b.Val }
func (b *Bool) Valid() bool { return b.Ok }

func (b *Bool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	(*b).Ok = asString != ""
	if asString == "true" {
		(*b).Val = true
	} else if asString == "false" {
		(*b).Val = false
	} else {
		return errors.New(
			fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString),
		)
	}

	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.Val {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}
