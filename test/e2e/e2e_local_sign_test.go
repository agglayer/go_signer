package e2e

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/signer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestLocalHash(t *testing.T) {
	ctx := context.TODO()
	sut, err := createLocalSigner(t, ctx, 1774)
	require.NoError(t, err)
	localSign, ok := sut.(*signer.LocalSign)
	require.True(t, ok)
	require.NoError(t, err)
	err = sut.Initialize(ctx)
	require.NoError(t, err)
	hash := crypto.Keccak256Hash([]byte("test"))
	signature, err := sut.SignHash(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, signature)
	// we have to remove the last byte V
	ok = localSign.Verify(hash, signature[0:64])
	require.True(t, ok)
}

func TestLocalSigner(t *testing.T) {
	testGenericSignerE2E(t, e2eTestParams{
		createSignerFunc: createLocalSigner,
		canSign:          true,
	})
}
