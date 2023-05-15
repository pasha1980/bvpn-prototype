package infrastructure

import (
	"bvpn-prototype/tests/testing"
	"os"
)

func I_have_sudo_permissions(t *testing.T) error {
	t.Assert(os.Getuid() == 0, "Non-root user")
	return nil
}
