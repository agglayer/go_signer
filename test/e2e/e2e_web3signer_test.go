package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	"github.com/agglayer/go_signer/test/e2e/helpers"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
)

// TestWeb3Signer tests the web3signer signer
func TestWeb3Signer(t *testing.T) {
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
	sign, err := signer.NewSigner(ctx, defaultChainID, signer.SignerConfig{
		Method: signer.MethodWeb3Signer,
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

	signed, err := sign.SignHash(ctx, common.Hash{})
	require.NoError(t, err)
	require.NotNil(t, signed)
	log.Debugf("signed hash: %s", common.Bytes2Hex(signed))
	require.Equal(t, "e82ed51b2b3964a6779171ee6589b1b2f5b5ebb77c1555626205d4619cb8df271a3f5c43f6b0ea3c76d852252d8a19539aa3ca2cb9fb66af3ac4dee7e846b4321c",
		common.Bytes2Hex(signed))

	testSendEthTx(t, sign.PublicAddress(), sign)

}
