package signer

import (
	"context"
	"fmt"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestMockSignRandomPrivateKey(t *testing.T) {
	sut, err := NewMockSign("test-mock-signer",
		log.WithFields("module", "test"),
		NewMockSignerConfig(""),
		1774)
	require.NoError(t, err, "should not return error when creating mock signer with random private key")
	ctx := context.TODO()
	signedHash, err := sut.SignHash(ctx, common.Hash{})
	require.NoError(t, err)
	err = sut.Verify(common.Hash{}, signedHash)
	require.NoError(t, err)
	fmt.Print("public addr: ", sut.PublicAddress())
}

func TestMockSignSpecificPrivateKey(t *testing.T) {
	sut, err := NewMockSign("test-mock-signer",
		log.WithFields("module", "test"),
		NewMockSignerConfig(testPrivateKeyHex),
		1774)
	require.NoError(t, err)
	ctx := context.TODO()
	signedHash, err := sut.SignHash(ctx, common.Hash{})
	require.NoError(t, err)
	err = sut.Verify(common.Hash{}, signedHash)
	require.NoError(t, err)
	require.Equal(t, testPublicKeyHex, sut.PublicAddress().Hex())
	fmt.Print("public addr: ", sut.PublicAddress())
}
