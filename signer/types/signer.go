package types

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNotImplemented       = fmt.Errorf("not implemented")
	ErrMissingConfigParam   = fmt.Errorf("missing config parameter")
	ErrBadConfigParams      = fmt.Errorf("bad config parameters")
	ErrNotImplementedMethod = fmt.Errorf("not implemented method")
)

type Signer interface {
	// Initialize the signer
	Initialize(context.Context) error

	// PublicKey returns the public key of the signer
	PublicAddress() common.Address
	// String returns a string representation of the signer (no secrets)
	String() string

	HashSigner
	TxSigner
}

type HashSigner interface {
	// SignHash signs the hash using the private key
	SignHash(context.Context, common.Hash) ([]byte, error)
}

type TxSigner interface {
	// SignTx signs the hash using the private key
	SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error)
}
