package infrastructure

import (
	"bvpn-prototype/tests/testing"
	"os"
)

func It_is_possible_to_create_file(t *testing.T) error {
	_, err := os.Create("test")
	t.Assert(err == nil, "Failed to create file")

	err = os.Remove("test")
	return err
}

func It_is_possible_to_create_dir(t *testing.T) error {
	err := os.Mkdir("test", 0666)
	t.Assert(err == nil, "Failed to create directory")

	err = os.Remove("test")
	return err
}
