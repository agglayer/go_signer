package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/stretchr/testify/require"
)

func createGCPSigner(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error) {
	t.Helper()
	KeyName := os.Getenv("GCP_KEY_NAME")
	require.NotEmpty(t, KeyName, "required env var GCP_KEY_NAME")
	log.Info("Creating GCP signer with key name: ", KeyName)
	return signer.NewSigner(ctx, chainID, signertypes.SignerConfig{
		Method: signertypes.MethodGCPKMS,
		Config: map[string]interface{}{
			"KeyName": KeyName,
		},
	}, "test", log.WithFields("module", "test"))
}

func TestGCPSigner(t *testing.T) {
	//t.Skip("It's not working yet")
	testGenericSignerE2E(t, e2eTestParams{
		createSignerFunc: createGCPSigner,
		canSign:          true,
	})
}
