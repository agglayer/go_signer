package signer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	signercommon "github.com/agglayer/go_signer/common"
	"github.com/agglayer/go_signer/log"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FieldPath     = "path"
	FieldPassword = "password"
)

var (
	ErrNoPrivateKey = fmt.Errorf("private key is nil")
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
		return fmt.Errorf("%s. Err: %w", e.logPrefix(), ErrNoPrivateKey)
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
		return nil, fmt.Errorf("%s. Err: %w", e.logPrefix(), ErrNoPrivateKey)
	}
	return crypto.Sign(hash.Bytes(), e.privateKey)
}

func (e *LocalSign) Verify(hash common.Hash, signature []byte) bool {
	if e.privateKey == nil {
		e.logger.Errorf("%s. Err: %w", e.logPrefix(), ErrNoPrivateKey)
		return false
	}
	pub := crypto.FromECDSAPub(&e.privateKey.PublicKey)
	log.Info("Pubkey: ", common.Bytes2Hex(pub))
	// If signature is longer than 64 bytes, we need to trim it. Usually it is 65 bytes
	// and the last byte is V (recovery id) that we don't need for verification.
	// because VerifySignature expects "signature should have the 64 byte [R || S] format."
	if len(signature) > 64 {
		signature = signature[0:64]
	}
	return crypto.VerifySignature(pub, hash.Bytes(), signature)
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
