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

const shutdownDockerAfterTest = false
const dockerIsAlreadyRunning = false

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

	sign, err := signer.NewSigner("test", log.WithFields("module", "test"), nil, signer.SignerConfig{
		Method: signer.MethodWeb3Signer,
		Config: map[string]interface{}{
			signer.FieldURL: "http://localhost:9999",
		},
	})
	require.NoError(t, err)
	ctx := context.TODO()
	err = sign.Initialize(ctx)
	require.NoError(t, err)
	expectedPublicAddress := "0x5497B14135C12dE8b230932e014921fEa0814421"
	require.Equal(t, expectedPublicAddress, sign.PublicAddress().String())

	signed, err := sign.SignHash(ctx, common.Hash{})
	require.NoError(t, err)
	require.NotNil(t, signed)
	log.Debugf("signed hash: %s", common.Bytes2Hex(signed))
	require.Equal(t, "08b11a644a17a08dbcb1f59d2d67cb6ee6bfebfcb97516b5627a51119ba057e147c30bb7730a2182ca87ce3fc295fdd64429a66c22ceae6f7c9c08a6f9516a111b",
		common.Bytes2Hex(signed))
}
