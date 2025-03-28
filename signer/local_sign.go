package signer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	signercommon "github.com/agglayer/go_signer/common"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FieldPath     = "path"
	FieldPassword = "password"
	// EIP155 is not enabled because the rest of user of this library
	// expected that signing a hash doesn't add anything else.
	enableEIP155 = false
)

// LocalSign is a signer that uses a local keystore file
type LocalSign struct {
	name          string
	logger        signercommon.Logger
	file          signercommon.KeystoreFileConfig
	privateKey    *ecdsa.PrivateKey
	publicAddress common.Address

	chainID uint64
	auth    *bind.TransactOpts
}

// NewLocalSignerConfig creates a generic config  (SignerConfig)
func NewLocalSignerConfig(path, pass string) signertypes.SignerConfig {
	return signertypes.SignerConfig{
		Method: signertypes.MethodLocal,
		Config: map[string]interface{}{
			FieldPath:     path,
			FieldPassword: pass,
		},
	}
}

// NewLocalConfig creates a KeystoreFileConfig (specific config) from a SignerConfig
func NewLocalConfig(cfg signertypes.SignerConfig) (signercommon.KeystoreFileConfig, error) {
	var res signercommon.KeystoreFileConfig
	// If there are no field in the config, return empty config
	// but if there are some field must match the expected ones
	if len(cfg.Config) == 0 {
		return signercommon.KeystoreFileConfig{}, nil
	}
	pathStr, ok := cfg.Config[FieldPath].(string)
	if !ok {
		return res, fmt.Errorf("field path is not string %v", cfg.Config[FieldPath])
	}
	passStr, ok := cfg.Config[FieldPassword].(string)
	if !ok {
		return res, fmt.Errorf("field pass is not string")
	}
	res = signercommon.KeystoreFileConfig{
		Path:     pathStr,
		Password: passStr,
	}
	return res, nil
}

// NewLocalSign creates a new LocalSign based on config
// name is the name of the signer
// logger is the logger to use
// file is the keystore file config
// chainID is the chainID to use (required to sync tx)
func NewLocalSign(name string, logger signercommon.Logger,
	file signercommon.KeystoreFileConfig, chainID uint64) *LocalSign {
	return &LocalSign{
		name:    name,
		logger:  logger,
		file:    file,
		chainID: chainID,
	}
}

// NewLocalSignFromPrivateKey creates a new LocalSign based on a private key
func NewLocalSignFromPrivateKey(name string,
	logger signercommon.Logger,
	privateKey *ecdsa.PrivateKey) *LocalSign {
	return &LocalSign{
		name:          name,
		logger:        logger,
		privateKey:    privateKey,
		publicAddress: crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

// Initialize initializes the LocalSign, read key if needed
func (e *LocalSign) Initialize(ctx context.Context) error {
	if err := e.initializeKey(); err != nil {
		return fmt.Errorf("%s failed to initialize key: %w", e.logPrefix(), err)
	}
	if err := e.initializeAuth(); err != nil {
		return fmt.Errorf("%s failed to initialize auth: %w", e.logPrefix(), err)
	}
	return nil
}

func (e *LocalSign) initializeKey() error {
	// Check if it's already initialized
	if e.privateKey != nil {
		return nil
	}
	privateKey, err := signercommon.NewKeyFromKeystore(e.file)
	if err != nil {
		return fmt.Errorf("%s initializeKey fails. Err: %w", e.logPrefix(), err)
	}
	if privateKey == nil {
		// If the private key is nil, the address is also nil
		// we allow to have a nil private key, it will fail if try to use it
		e.logger.Warnf("%s private key is nil", e.logPrefix())
		return nil
	}
	e.privateKey = privateKey
	e.publicAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
	return nil
}

func (e *LocalSign) initializeAuth() error {
	if e.auth != nil {
		return nil
	}
	if e.privateKey == nil {
		return nil
	}
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, new(big.Int).SetUint64(e.chainID))
	if err != nil {
		return fmt.Errorf("%s can't initialize auth. Err: %w", e.logPrefix(), err)
	}
	e.auth = auth
	return nil
}

// SignHash signs a hash
func (e *LocalSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	if e.privateKey == nil {
		return nil, fmt.Errorf("%s private key is nil", e.logPrefix())
	}
	if enableEIP155 {
		// length of the hash is 32 bytes, so it's hardcoded
		hashWithPrefix := crypto.Keccak256(append([]byte("\x19Ethereum Signed Message:\n32"), hash.Bytes()...))
		sig, err := crypto.Sign(hashWithPrefix, e.privateKey)
		if err != nil {
			return nil, fmt.Errorf("%s can't sign hash. Err: %w", e.logPrefix(), err)
		}
		// Set r(recoveryID) as eth_sign that is 27 or 28
		// crypto.Sign returns  0 or 1
		sig[64] += 27
		return sig, err
	}
	return crypto.Sign(hash.Bytes(), e.privateKey)
}

func (e *LocalSign) PublicAddress() common.Address {
	return e.publicAddress
}

func (e *LocalSign) String() string {
	return fmt.Sprintf("%s path:%s, pubAddr: %s", e.logPrefix(), e.file, e.publicAddress.String())
}

func (e *LocalSign) logPrefix() string {
	return fmt.Sprintf("signer: %s[%s]: ", signertypes.MethodLocal, e.name)
}

func (e *LocalSign) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	if e.auth == nil {
		return nil, fmt.Errorf("%s can't signTx because auth is nil", e.logPrefix())
	}

	signedTx, err := e.auth.Signer(e.publicAddress, tx)
	if err != nil {
		return nil, fmt.Errorf("%s can't signTx because auth.Signer returns error %w", e.logPrefix(), err)
	}
	return signedTx, nil
}
