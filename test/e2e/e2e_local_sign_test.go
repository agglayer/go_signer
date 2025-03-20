package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	"github.com/agglayer/go_signer/test/e2e/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestLocalSigner(t *testing.T) {
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

	password, err := os.ReadFile("key_store/funded_addr.password")
	require.NoError(t, err)
	trimmedPassword := strings.TrimSpace(string(password))
	sign, err := signer.NewSigner(ctx, chainID.Uint64(), signer.SignerConfig{
		Method: signer.MethodLocal,
		Config: map[string]interface{}{
			signer.FieldPath:     "key_store/funded_addr",
			signer.FieldPassword: trimmedPassword,
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
