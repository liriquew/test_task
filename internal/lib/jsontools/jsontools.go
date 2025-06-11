package jsontools

import (
	"io"

	"encoding/json"
)

func Encode(w io.Writer, obj any) error {
	return json.NewEncoder(w).Encode(obj)
}

func Decode(r io.Reader, obj any) error {
	return json.NewDecoder(r).Decode(obj)
}
