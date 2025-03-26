package helpers

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const bash = "/bin/bash"

func runCommand(t *testing.T, mustBeOk bool, command string) {
	t.Helper()
	msg, err := exec.Command(bash, "-l", "-c", command).CombinedOutput()
	if err != nil {
		t.Log(err)
	}
	if mustBeOk {
		require.NoError(t, err, string(msg))
	}
}
