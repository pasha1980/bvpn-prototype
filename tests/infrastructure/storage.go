package infrastructure

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/tests/testing"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
)

func All_files_readable(t *testing.T) error {
	requiredFiles := []string{
		"bvpn.log",
		"bvpn-errors.log",
		"peers.bvpn",
		"prv.pem",
		"chain.db",
	}

	for _, file := range requiredFiles {
		_, err := os.ReadFile(file)
		t.Assert(err == nil || os.IsNotExist(err), "Filed to read file "+file)
	}
	return nil
}

func It_is_possible_to_initialize_database(t *testing.T) error {
	db, err := sql.Open("sqlite3", "test.db")
	t.Assert(err == nil, "Failed to initialize db")

	_, err = db.Exec("select 1")
	t.Assert(err == nil, "Failed to communicate with db")

	err = os.Remove("test.db")
	return err
}

func It_is_correct_mempool_format(t *testing.T) error {
	base64Encoded, err := os.ReadFile("mempool.bvpn")
	if os.IsNotExist(err) {
		return nil
	}
	t.Assert(err == nil, "Failed to read mempool")

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	t.Assert(err == nil, "Failed to decode mempool")

	var storage map[string]block_data.ChainStored
	err = json.Unmarshal(jsonEncoded, &storage)
	t.Assert(err == nil, "Failed to decode mempool")

	return nil
}

func It_is_correct_peer_storage_format(t *testing.T) error {
	base64Encoded, err := os.ReadFile("peers.bvpn")
	if os.IsNotExist(err) {
		return nil
	}
	t.Assert(err == nil, "Failed to read peers storage")

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	t.Assert(err == nil, "Failed to decode peers storage")

	var storage struct {
		Data []entity.Node
	}
	err = json.Unmarshal(jsonEncoded, &storage)
	t.Assert(err == nil, "Failed to decode peer storage")
	return nil
}

func It_is_correct_prv_key_format(t *testing.T) error {
	data, err := os.ReadFile("prv.pem")
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}
	block, _ := pem.Decode(data)
	t.Assert(block != nil, "Failed to parse prv.pem")
	return nil
}

func It_is_correct_profile_storage_format(t *testing.T) error {
	return nil // todo
}
