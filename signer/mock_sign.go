package signer

import (
	"context"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	goethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// MockSign is a signer that uses an arbiratrary private key for testing purposes.
// basically it's a wrapper over LocalSign that, instead of getting the private key from
// a keystore file, it uses a private key that is set in the configuration.
type MockSign struct {
	name      string
	logger    signercommon.Logger
	cfg       MockSignConfigure
	localSign *LocalSign
}

func NewMockSign(name string, logger signercommon.Logger,
	genericCfg signertypes.SignerConfig, chainID uint64) (*MockSign, error) {
	cfg, err := NewMockConfig(genericCfg)
	if err != nil {
		return nil, err
	}

	privateKey := cfg.privateKey
	if privateKey == nil {
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate private key: %w", err)
		}
	}
	return &MockSign{
		name:      name,
		logger:    logger,
		cfg:       cfg,
		localSign: NewLocalSignFromPrivateKey("MockSign("+name+")", logger, privateKey, chainID),
	}, nil
}

func (e *MockSign) String() string {
	return fmt.Sprintf("MockSign{name:%s, mode:%s}", e.name, e.cfg.mode.String())
}
func (e *MockSign) Initialize(ctx context.Context) error {
	// It needs to initialize bin.Auth
	return e.localSign.Initialize(ctx)
}

// SignHash signs a hash
func (e *MockSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return e.localSign.SignHash(ctx, hash)
}

func (e *MockSign) Verify(hash common.Hash, signature []byte) error {
	return e.localSign.Verify(hash, signature)
}

func (e *MockSign) PublicAddress() common.Address {
	return e.localSign.PublicAddress()
}
func (e *MockSign) SignTx(ctx context.Context, tx *goethereumtypes.Transaction) (*goethereumtypes.Transaction, error) {
	return e.localSign.SignTx(ctx, tx)
}
