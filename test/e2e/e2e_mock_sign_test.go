package e2e

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
)

func createMockSigner(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error) {
	t.Helper()
	return signer.NewMockSign("test-mock-signer",
		log.WithFields("module", "test"),
		signer.NewMockSignerConfig("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"),
		chainID)
}

func TestMockSigner(t *testing.T) {
	testGenericSignerE2E(t, e2eTestParams{
		createSignerFunc: createMockSigner,
		canSign:          true,
	})
}
