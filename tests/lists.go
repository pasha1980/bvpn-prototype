package tests

import (
	"bvpn-prototype/tests/infrastructure"
	"bvpn-prototype/tests/testing"
)

var infrastructureTestList = []func(t *testing.T) error{
	infrastructure.It_is_possible_to_create_file,
	infrastructure.It_is_possible_to_create_dir,
	infrastructure.All_files_readable,
	infrastructure.It_is_possible_to_initialize_database,
	infrastructure.It_is_correct_mempool_format,
	infrastructure.It_is_correct_peer_storage_format,
	infrastructure.It_is_correct_profile_storage_format,
	infrastructure.It_is_correct_profile_storage_format,
	infrastructure.I_have_sudo_permissions,
	infrastructure.It_is_possible_to_use_ip_command,
	infrastructure.It_is_stable_internet_conenction,
	infrastructure.It_is_correct_config_file,
}

var featureTestsList = []func(t *testing.T) error{}

var unitTestsList = []func(t *testing.T) error{}
