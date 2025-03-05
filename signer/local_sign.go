package signer

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FieldPath     = "path"
	FieldPassword = "password"
)

// LocalSign is a signer that uses a local keystore file
type LocalSign struct {
	name          string
	logger        signercommon.Logger
	file          signercommon.KeystoreFileConfig
	privateKey    *ecdsa.PrivateKey
	publicAddress common.Address
}

// NewLocalSignerConfig creates a generic config  (SignerConfig)
func NewLocalSignerConfig(path, pass string) SignerConfig {
	return SignerConfig{
		Method: MethodLocal,
		Config: map[string]interface{}{
			FieldPath:     path,
			FieldPassword: pass,
		},
	}
}

// NewLocalConfig creates a KeystoreFileConfig (specific config) from a SignerConfig
func NewLocalConfig(cfg SignerConfig) (signercommon.KeystoreFileConfig, error) {
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
func NewLocalSign(name string, logger signercommon.Logger, file signercommon.KeystoreFileConfig) *LocalSign {
	return &LocalSign{
		name:   name,
		logger: logger,
		file:   file,
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
	// Check if it's already initialized
	if e.privateKey != nil {
		return nil
	}
	privateKey, err := signercommon.NewKeyFromKeystore(e.file)
	if err != nil {
		return err
	}
	if privateKey == nil {
		// If the private key is nil, the address is also nil
		return nil
	}
	e.privateKey = privateKey
	e.publicAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
	return nil
}

// SignHash signs a hash
func (e *LocalSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	if e.privateKey == nil {
		return nil, fmt.Errorf("%s private key is nil", e.logPrefix())
	}
	return crypto.Sign(hash.Bytes(), e.privateKey)
}

func (e *LocalSign) PublicAddress() common.Address {
	return e.publicAddress
}

func (e *LocalSign) String() string {
	return fmt.Sprintf("signer: %s path:%s, pubAddr: %s", e.logPrefix(), e.file, e.publicAddress.String())
}

func (e *LocalSign) logPrefix() string {
	return fmt.Sprintf("signer:%s[%s]: ", MethodLocal, e.name)
}
