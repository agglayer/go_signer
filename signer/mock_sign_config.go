package signer

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FieldMockPrivateKey = "privatekey"
	FieldMockPublicKey  = "publickey"
)

type MockSignMode int

var (
	MockSignModeRandom     MockSignMode = 0 // MockSignModeRandom it generate a random private key
	MockSignModePrivateKey MockSignMode = 1 // MockSignModePublicKey have set a specific private key
	MockSignModePublicKey  MockSignMode = 2 // MockSignModePublicKey no crypt but it respond with a public key
)

func (s MockSignMode) String() string {
	return []string{
		"Random",
		"PrivateKey",
		"NoCrypt+PublicKey",
	}[int(s)]
}

type MockSignConfigure struct {
	privateKey    *ecdsa.PrivateKey
	publicAddress *common.Address
	mode          MockSignMode
}

func (c *MockSignConfigure) String() string {
	if c == nil {
		return "MockSignConfigure{nil}"
	}
	res := fmt.Sprintf("MockSignConfigure{mode: %s", c.mode.String())
	if c.privateKey != nil {
		res += ", privateKey: SET"
	}
	if c.publicAddress != nil {
		res += fmt.Sprintf(", publicAddress: %s", c.publicAddress.Hex())
	}
	return res + "}"
}

// NewMockSignerConfig creates a new mock signer configuration (the general one that is on configfiles)
func NewMockSignerConfig(privateKey, publicKey string) signertypes.SignerConfig {
	res := signertypes.SignerConfig{
		Method: signertypes.MethodMock,
		Config: map[string]interface{}{},
	}
	if privateKey != "" {
		res.Config[FieldMockPrivateKey] = privateKey
	}
	if publicKey != "" {
		res.Config[FieldMockPublicKey] = publicKey
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

	publicKeyStr, err := cfg.Get(FieldMockPublicKey)
	if err != nil && !errors.Is(err, signertypes.ErrMissingConfigParam) {
		return res, fmt.Errorf("config %s: error in field %s . Err: %w",
			cfg.Method, FieldMockPublicKey, err)
	}

	if privateKeyStr != "" {
		err := res.LoadPrivateKey(privateKeyStr)
		if err != nil {
			return res, fmt.Errorf("config %s: field %s is not a valid private key. Err: %w", cfg.Method, FieldMockPrivateKey, err)
		}
		res.mode = MockSignModePrivateKey
	}

	if publicKeyStr != "" {
		err := res.SetPublicKey(publicKeyStr)
		if err != nil {
			return res, fmt.Errorf("config %s: field %s is not a valid public key: %w", cfg.Method, FieldMockPublicKey, err)
		}
		if res.privateKey != nil {
			res.mode = MockSignModePublicKey
		}
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
	//publicAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)
	//c.publicAddress = &publicAddress
	return nil
}

func (c *MockSignConfigure) SetPublicKey(hexKey string) error {
	hexKey = remove0xPrefix(hexKey) // Remove '0x' prefix if present
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key hex: %w", err)
	}
	if len(key) != common.AddressLength {
		return fmt.Errorf("invalid public key length: expected %d bytes, got %d", common.AddressLength, len(key))
	}
	pubKey := common.BytesToAddress(key)
	if c.privateKey != nil {
		// If there are a private key set, we only accept the public key that matches it
		publicAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)
		if publicAddress != pubKey {
			return fmt.Errorf("public key does not match existing public address: %s != %s", publicAddress.Hex(), pubKey.Hex())
		}
		// We don't set c.publicAddress because is not required in this case
		return nil
	}
	c.publicAddress = &pubKey
	return nil
}
