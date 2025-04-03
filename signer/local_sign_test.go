package signer

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	signercommon "github.com/agglayer/go_signer/common"
	"github.com/agglayer/go_signer/log"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// To keep compatibility with previous version, an empty config file
// meant there is no privateKey (nil), so the idea is keep the same behavior
func TestNewKeyStoreFileConfigEmpty(t *testing.T) {
	cfg, err := NewLocalConfig(signertypes.SignerConfig{})
	require.NoError(t, err)
	require.Equal(t, "", cfg.Path)
	require.Equal(t, "", cfg.Password)
}

func TestNewLocalSignerConfig(t *testing.T) {
	cfg := NewLocalSignerConfig("/app/sequencer.keystore", "test")
	require.Equal(t, signertypes.MethodLocal, cfg.Method)
	require.Equal(t, "/app/sequencer.keystore", cfg.Config[FieldPath])
	require.Equal(t, "test", cfg.Config[FieldPassword])

	localCfg, err := NewLocalConfig(cfg)
	require.NoError(t, err)
	require.Equal(t, "/app/sequencer.keystore", localCfg.Path)
	require.Equal(t, "test", localCfg.Password)
}

func TestNewLocalSignerConfigWrongData(t *testing.T) {
	cfg := NewLocalSignerConfig("/app/sequencer.keystore", "test")
	cfg.Config[FieldPath] = 123
	_, err := NewLocalConfig(cfg)
	require.Error(t, err)
}

func TestNewLocalSign(t *testing.T) {
	sut := NewLocalSign("name", nil, signercommon.KeystoreFileConfig{}, 0)
	require.NotNil(t, sut)
	require.Equal(t, "name", sut.name)
	require.Nil(t, sut.logger)
	require.Equal(t, "", sut.file.Path)
}

func TestNewLocalSignFromPrivateKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	sut := NewLocalSignFromPrivateKey("name", nil, privateKey)
	require.NotNil(t, sut)
	ctx := context.TODO()
	err = sut.Initialize(ctx)
	require.NoError(t, err)
	pubAddr := sut.PublicAddress()
	require.NotNil(t, pubAddr)
	t.Log("pubAddr: ", pubAddr.String())
	str := sut.String()
	require.NotEmpty(t, str)
	hashToSign := common.Hash{}
	signature, err := sut.SignHash(ctx, hashToSign)
	require.NoError(t, err)
	signOk := sut.Verify(hashToSign, signature[0:64])
	require.True(t, signOk)
}

func TestNewLocalSignEmpty(t *testing.T) {
	logger := log.WithFields("test", "test")
	sut := NewLocalSign("name", logger, signercommon.KeystoreFileConfig{}, 0)
	err := sut.Initialize(context.Background())
	require.NoError(t, err)
	pubAddr := sut.PublicAddress()
	require.Equal(t, common.Address{}, pubAddr)
	_, err = sut.SignHash(context.Background(), common.Hash{})
	require.Error(t, err)
	require.NotEmpty(t, sut.String())
}
