package infrastructure

import (
	"bvpn-prototype/tests/testing"
	"gopkg.in/yaml.v3"
	"os"
)

func It_is_correct_config_file(t *testing.T) error {
	configFile := "./config.yaml"
	file, err := os.ReadFile(configFile)
	t.Assert(err == nil, "Failed to open config file "+configFile)
	if err != nil {
		return err
	}

	testMap := make(map[string]any)
	err = yaml.Unmarshal(file, &testMap)
	t.Assert(err == nil, "Failed to parse config file")
	return err
}
