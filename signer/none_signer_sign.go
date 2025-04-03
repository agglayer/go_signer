package signer

import (
	"context"

	gosignertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NoneSign is a signer that does not sign anything
type NoneSign struct {
}

// Initialize initializes the NoneSign signer
func (s *NoneSign) Initialize(ctx context.Context) error {
	return nil
}

// PublicAddress returns 0x0 address
func (s *NoneSign) PublicAddress() common.Address {
	return common.Address{}
}

// String returns the string representation of the NoneSign signer
func (s *NoneSign) String() string {
	return "none"
}

// SignHash returns error always
func (s *NoneSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return nil, gosignertypes.ErrNotImplemented
}

// SignTx returns error always
func (s *NoneSign) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	return nil, gosignertypes.ErrNotImplemented
}
