package helpers

import (
	"testing"
)

func StartWeb3SignerDocker(t *testing.T) {
	t.Helper()
	runCommand(t, true, "docker run -d -p 9999:9000 --name web3signer consensys/web3signer:develop --http-listen-port=9000 eth1 --chain-id 123")
}

func StopWeb3SignerDocker(t *testing.T) {
	t.Helper()
	runCommand(t, false, "docker stop web3signer; docker rm web3signer")
}
