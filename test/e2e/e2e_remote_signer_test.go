package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/agglayer/go_signer/test/e2e/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TestRemoteSigner tests the web3signer signer
func TestRemoteSigner(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	if !dockerIsAlreadyRunning {
		dockerCompose := helpers.NewDockerCompose()
		dockerCompose.Down(t)
		dockerCompose.Up(t)
		defer func() {
			if shutdownDockerAfterTest {
				dockerCompose.Down(t)
			}
		}()
		dockerCompose.WaitHealthy(t, 40*time.Second)
	}
	ctx := context.TODO()
	sign, err := signer.NewSigner(ctx, defaultChainID, signertypes.SignerConfig{
		Method: signertypes.MethodRemoteSigner,
		Config: map[string]interface{}{
			signer.FieldURL:     "http://localhost:9999",
			signer.FieldAddress: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
	}, "test", log.WithFields("module", "test"))
	require.NoError(t, err)

	err = sign.Initialize(ctx)
	require.NoError(t, err)
	expectedPublicAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	require.Equal(t, expectedPublicAddress, sign.PublicAddress().String())

	// Can't sign using eth_sign because EIP155 prefix the data with a string
	// and the aggkit expect to have a signature of just the hash
	_, err = sign.SignHash(ctx, common.Hash{})
	require.Error(t, err)

	testSendEthTx(t, sign.PublicAddress(), sign)
}
