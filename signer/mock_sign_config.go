package signer

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FieldMockPrivateKey = "privatekey"
)

type MockSignMode int

var (
	MockSignModeRandom     MockSignMode = 0 // MockSignModeRandom it generate a random private key
	MockSignModePrivateKey MockSignMode = 1 // MockSignModePublicKey have set a specific private key
)

func (s MockSignMode) String() string {
	return []string{
		"Random",
		"PrivateKey",
	}[int(s)]
}

type MockSignConfigure struct {
	privateKey *ecdsa.PrivateKey
	mode       MockSignMode
}

func (c *MockSignConfigure) String() string {
	if c == nil {
		return "MockSignConfigure{nil}"
	}
	res := fmt.Sprintf("MockSignConfigure{mode: %s", c.mode.String())
	if c.privateKey != nil {
		res += ", privateKey: SET (public:" + crypto.PubkeyToAddress(c.privateKey.PublicKey).Hex() + ")"
	}
	return res + "}"
}

// NewMockSignerConfig creates a new mock signer configuration (the general one that is on configfiles)
func NewMockSignerConfig(privateKey string) signertypes.SignerConfig {
	res := signertypes.SignerConfig{
		Method: signertypes.MethodMock,
		Config: map[string]interface{}{},
	}
	if privateKey != "" {
		res.Config[FieldMockPrivateKey] = privateKey
	}
	return res
}

// NewMockConfig creates a MockSignConfigure specific for this signer
func NewMockConfig(cfg signertypes.SignerConfig) (MockSignConfigure, error) {
	res := MockSignConfigure{
		mode: MockSignModeRandom, // Default mode is random
	}
	// If there are no field in the config, return empty config
	if len(cfg.Config) == 0 {
		return MockSignConfigure{}, nil
	}
	privateKeyStr, err := cfg.Get(FieldMockPrivateKey)
	if err != nil && !errors.Is(err, signertypes.ErrMissingConfigParam) {
		return res, fmt.Errorf("config %s: error in field %s . Err: %w",
			cfg.Method, FieldMockPrivateKey, err)
	}

	if privateKeyStr != "" {
		err := res.LoadPrivateKey(privateKeyStr)
		if err != nil {
			return res, fmt.Errorf("config %s: field %s is not a valid private key. Err: %w",
				cfg.Method, FieldMockPrivateKey, err)
		}
		res.mode = MockSignModePrivateKey
	}
	return res, nil
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func remove0xPrefix(str string) string {
	if has0xPrefix(str) {
		return str[2:] // Remove '0x' prefix if present
	}
	return str
}

func (c *MockSignConfigure) LoadPrivateKey(hexKey string) error {
	hexKey = remove0xPrefix(hexKey) // Remove '0x' prefix if present
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return fmt.Errorf("failed to decode private key hex: %w", err)
	}
	c.privateKey, err = crypto.ToECDSA(key)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	return nil
}
