package helpers

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func startWeb3SignerDocker(t *testing.T) {
	t.Helper()
	stopWeb3SignerDocker(t)
	msg, err := exec.Command("bash", "-l", "-c", "docker run -d -p 9999:9000 --name web3signer consensys/web3signer:develop --http-listen-port=9000 eth1 --chain-id 123").CombinedOutput()
	require.NoError(t, err, string(msg))
}

func stopWeb3SignerDocker(t *testing.T) {
	t.Helper()
	_, err := exec.Command("bash", "-l", "-c", "docker stop web3signer; docker rm web3signer").CombinedOutput()
	if err != nil {
		t.Log(err)
	}
}
