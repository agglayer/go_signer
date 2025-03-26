package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/agglayer/go_signer/test/e2e/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestGCPSigner(t *testing.T) {
	t.Skip("It's not working yet")
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
	ethClient, err := ethclient.Dial(gethURL)
	require.NoError(t, err)
	defer ethClient.Close()
	chainID, err := ethClient.ChainID(ctx)
	require.NoError(t, err)
	log.Info("chainID: ", chainID.Uint64())
	KeyName := os.Getenv("GCP_KEY_NAME")
	require.NotEmpty(t, KeyName, "required env var GCP_KEY_NAME")
	sign, err := signer.NewSigner(ctx, chainID.Uint64(), signertypes.SignerConfig{
		Method: signertypes.MethodGCPKMS,
		Config: map[string]interface{}{
			"KeyName": KeyName,
		},
	}, "test", log.WithFields("module", "test"))
	require.NoError(t, err)

	err = sign.Initialize(ctx)
	require.NoError(t, err)
	require.Equal(t, publicAddressTest, sign.PublicAddress().String(), "public_addr")

	signed, err := sign.SignHash(ctx, common.Hash{})
	require.NoError(t, err)
	require.NotNil(t, signed)
	log.Debugf("signed hash: %s", common.Bytes2Hex(signed))
	require.Equal(t, "b8823364c90ea0d2700d5ad0fe39d16778bc07ce7df4779ff35e4b2660d043cb74a002439225d1d518f9f1cf3db005f5e143196543fd5146a34bf63f0b810ade00",
		common.Bytes2Hex(signed))

	testSendEthTx(t, sign.PublicAddress(), sign)
}
