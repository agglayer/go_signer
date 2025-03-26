package signer

import (
	"context"
	"fmt"
	"testing"

	"github.com/agglayer/go_signer/log"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSignerExploratory(t *testing.T) {
	t.Skip("This test is for exploratory purposes only")
	logger := log.WithFields("test", "test")
	ctx := context.TODO()
	localConfig := signertypes.SignerConfig{
		Method: signertypes.MethodLocal,
		Config: map[string]interface{}{
			FieldPath:     "../tmp/local_config/sequencer.keystore",
			FieldPassword: "pSnv6Dh5s9ahuzGzH9RoCDrKAMddaX3m",
		},
	}

	local, err := NewSigner(ctx, 0, localConfig, "test-local", logger)
	require.NoError(t, err)
	require.NotNil(t, local)
	require.NoError(t, local.Initialize(ctx))
	require.NotNil(t, local.PublicAddress())
	hash := common.Hash{}
	sign, err := local.SignHash(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, sign)
	require.Equal(t, 65, len(sign))
	fmt.Print(sign)
	w3sConfig := signertypes.SignerConfig{
		Method: signertypes.MethodRemoteSigner,
		Config: map[string]interface{}{
			FieldURL: "http://localhost:9000",
		},
	}
	remoteSigner, err := NewSigner(ctx, 0, w3sConfig, "test-w3s", logger)
	require.NoError(t, err)
	require.NotNil(t, remoteSigner)
	require.NoError(t, remoteSigner.Initialize(ctx))
	require.NotNil(t, remoteSigner.PublicAddress())
	signW3s, err := remoteSigner.SignHash(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, signW3s)
	require.Equal(t, 65, len(signW3s))
	fmt.Print(signW3s)
}

func TestNewSigner(t *testing.T) {
	logger := log.WithFields("test", "test")
	ctx := context.TODO()
	t.Run("unknown signer method", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{Method: "unknown_method"}, "test", logger)
		require.Error(t, err)
		require.Nil(t, sut)
	})
	t.Run("empty method is local", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{}, "test", logger)
		require.NoError(t, err)
		require.NotNil(t, sut)
		require.Contains(t, sut.String(), signertypes.MethodLocal)
	})

	t.Run("wrong local config", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Config: map[string]interface{}{
				FieldPath: 1234,
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Nil(t, sut)
	})

	t.Run("wrong local config", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Config: map[string]interface{}{
				FieldPath: 1234,
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Nil(t, sut)
	})

	t.Run("wrong remote config", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Method: signertypes.MethodRemoteSigner,
			Config: map[string]interface{}{
				FieldAddress: 1234,
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Nil(t, sut)
	})

	t.Run("wrong remote config2", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Method: signertypes.MethodRemoteSigner,
			Config: map[string]interface{}{
				FieldAddress: "NOTHEXA",
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Contains(t, err.Error(), FieldAddress)
		require.Nil(t, sut)
	})

	t.Run("wrong remote config3", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Method: signertypes.MethodRemoteSigner,
			Config: map[string]interface{}{
				FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
				FieldURL:     1234,
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Contains(t, err.Error(), FieldURL)
		require.Nil(t, sut)
	})

	t.Run("wrong remote missing URL", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Method: signertypes.MethodRemoteSigner,
			Config: map[string]interface{}{
				FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
			},
		}, "test-local", logger)
		require.Error(t, err)
		require.Contains(t, err.Error(), FieldURL)
		require.Nil(t, sut)
	})

	t.Run("remote config", func(t *testing.T) {
		sut, err := NewSigner(ctx, 1, signertypes.SignerConfig{
			Method: signertypes.MethodRemoteSigner,
			Config: map[string]interface{}{
				FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
				FieldURL:     "http://localhost:9001",
			},
		}, "test-local", logger)
		require.NoError(t, err)
		require.NotNil(t, sut)
	})
}
