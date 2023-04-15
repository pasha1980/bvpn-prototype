package utils

import (
	"bytes"
	"encoding/gob"
)

func SizeOf(v interface{}) int {
	b := new(bytes.Buffer)
	_ = gob.NewEncoder(b).Encode(v)
	return b.Len()
}
