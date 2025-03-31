package signer

import (
	"context"

	gosignertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type NoneSign struct {
}

func (s *NoneSign) Initialize(ctx context.Context) error {
	return nil
}

func (s *NoneSign) PublicAddress() common.Address {
	return common.Address{}
}

func (s *NoneSign) String() string {
	return "none"
}

func (s *NoneSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return nil, gosignertypes.ErrNotImplemented
}

func (s *NoneSign) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	return nil, gosignertypes.ErrNotImplemented
}
