package helpers

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func startWeb3SignerDocker(t *testing.T) {
	t.Helper()
	runCommand(t, true, "docker run -d -p 9999:9000 --name web3signer consensys/web3signer:develop --http-listen-port=9000 eth1 --chain-id 123")
}

func stopWeb3SignerDocker(t *testing.T) {
	t.Helper()
	runCommand(t, false, "docker stop web3signer; docker rm web3signer")
}

func runCommand(t *testing.T, mustBeOk bool, command string) {
	t.Helper()
	msg, err := exec.Command("bash", "-l", "-c", command).CombinedOutput()
	if err != nil {
		t.Log(err)
	}
	if mustBeOk {
		require.NoError(t, err, string(msg))
	}
}
