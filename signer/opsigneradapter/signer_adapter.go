package opsigneradapter

import (
	"context"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	gosignertypes "github.com/agglayer/go_signer/signer/types"
	opsignerprovider "github.com/ethereum-optimism/infra/op-signer/provider"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type SignerAdapter struct {
	opSigner opsignerprovider.SignatureProvider
	ctx      context.Context
	keyName  string
}

var _ gosignertypes.Signer = (*SignerAdapter)(nil)

func NewSignerAdapter(ctx context.Context, opSigner opsignerprovider.SignatureProvider, keyName string) *SignerAdapter {
	return &SignerAdapter{
		opSigner: opSigner,
		ctx:      ctx,
		keyName:  keyName,
	}
}

func NewSignerAdapterFromConfig(ctx context.Context, logger signercommon.Logger,
	cfg gosignertypes.SignerConfig) (*SignerAdapter, error) {
	opConfig := opsignerprovider.ProviderConfig{
		ProviderType: opsignerprovider.ProviderType(cfg.Method),
	}
	opSigner, err := opsignerprovider.NewSignatureProvider(NewLoggerAdapter(logger), opConfig.ProviderType, opConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating opSignerProvider. Err: %w", err)
	}
	keyName, err := cfg.Get("KeyName")
	if err != nil {
		return nil, fmt.Errorf("error getting keyName from config. Err: %w", err)
	}
	return NewSignerAdapter(ctx, opSigner, keyName), nil
}

func (s *SignerAdapter) Initialize(context.Context) error {
	_, err := s.opSigner.GetPublicKey(s.ctx, s.keyName)
	if err != nil {
		return fmt.Errorf("fails to Initialize because error getting public key from opSigner. Err: %w", err)
	}
	return nil
}

func (s *SignerAdapter) PublicAddress() common.Address {
	res, err := s.opSigner.GetPublicKey(s.ctx, s.keyName)
	if err != nil {
		return common.Address{}
	}
	return convertPublicKeyToAddress(res)
}

func convertPublicKeyToAddress(publicKey []byte) common.Address {
	addr := crypto.Keccak256(publicKey[1:])[12:]
	return common.BytesToAddress(addr)
}

func (s *SignerAdapter) String() string {
	return "op_signer_adapter"
}

func (s *SignerAdapter) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return s.opSigner.SignDigest(ctx, s.keyName, hash[:])
}

func (s *SignerAdapter) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	//txSigner := types.LatestSignerForChainID(tx.ChainId())
	txSigner := types.HomesteadSigner{}
	digest := txSigner.Hash(tx)
	signature, err := s.opSigner.SignDigest(ctx, s.keyName, digest.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error signTx opSigner.SignDigest. Err: %w ", err)
	}
	signed, err := tx.WithSignature(txSigner, signature)
	if err != nil {
		return nil, fmt.Errorf("error signTx tx.WithSignature. Err: %w ", err)
	}
	return signed, nil
}
