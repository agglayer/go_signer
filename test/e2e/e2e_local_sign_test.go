package e2e

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/stretchr/testify/require"
)

func createLocalSigner(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error) {
	password, err := os.ReadFile("key_store/funded_addr.password")
	require.NoError(t, err)
	trimmedPassword := strings.TrimSpace(string(password))
	return signer.NewSigner(ctx, chainID, signertypes.SignerConfig{
		Method: signertypes.MethodLocal,
		Config: map[string]interface{}{
			signer.FieldPath:     "key_store/funded_addr",
			signer.FieldPassword: trimmedPassword,
		},
	}, "test", log.WithFields("module", "test"))
}
func TestLocalSigner(t *testing.T) {
	testGenericSignerE2E(t, createLocalSigner)
}
