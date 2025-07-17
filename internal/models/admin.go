package models

import (
	"database/sql"
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
			sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
	}
}

type Bool struct {
	sql.NullBool
}

func (b *Bool) Valid() bool { return b.NullBool.Bool }

func (b *Bool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	(*b).NullBool.Valid = asString != ""
	if asString == "true" {
		(*b).NullBool.Bool = true
	} else if asString == "false" {
		(*b).NullBool.Bool = false
	} else {
		return fmt.Errorf(
			"Boolean unmarshal error: invalid input %s", asString,
		)
	}

	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.NullBool.Bool {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}
