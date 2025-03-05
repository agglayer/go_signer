package common

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// KeystoreFileConfig has all the information needed to load a private key from a key store file
type KeystoreFileConfig struct {
	// Path is the file path for the key store file
	Path string `mapstructure:"Path"`

	// Password is the password to decrypt the key store file
	Password string `mapstructure:"Password"`
}

// NewKeyFromKeystore creates a private key from a keystore file
func NewKeyFromKeystore(cfg KeystoreFileConfig) (*ecdsa.PrivateKey, error) {
	if cfg.Path == "" && cfg.Password == "" {
		return nil, nil
	}
	keystoreEncrypted, err := os.ReadFile(filepath.Clean(cfg.Path))
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keystoreEncrypted, cfg.Password)
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}
