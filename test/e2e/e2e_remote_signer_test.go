package e2e

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
)

func createRemoteSigner(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error) {
	t.Helper()
	return signer.NewSigner(ctx, defaultChainID, signertypes.SignerConfig{
		Method: signertypes.MethodRemoteSigner,
		Config: map[string]interface{}{
			signer.FieldURL:     "http://localhost:9999",
			signer.FieldAddress: publicAddressTest,
		},
	}, "test", log.WithFields("module", "test"))
}

// TestRemoteSigner tests the web3signer signer
func TestRemoteSigner(t *testing.T) {
	testGenericSignerE2E(t, createRemoteSigner)
}
