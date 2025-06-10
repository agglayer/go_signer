package signer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	privateKey = "0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0"
	publicKey  = "0xc653eCD4AC5153a3700Fb13442Bcf00A691cca16"
	publicKey2 = "0x1234567890abcdef1234567890abcdef12345678"
)

func TestGenereteSpecificOnlyPrivateKey(t *testing.T) {
	genericCfg := NewMockSignerConfig(privateKey, "")
	specificCfg, err := NewMockConfig(genericCfg)
	require.NoError(t, err, "should not return error when creating specific config from generic config")
	require.NotNil(t, specificCfg, "specific config should not be nil")
	require.NotNil(t, specificCfg.privateKey, "specific config private key should not be nil")
	require.Nil(t, specificCfg.publicAddress)
	require.Equal(t, "MockSignConfigure{mode: PrivateKey, privateKey: SET}", specificCfg.String())
}

func TestGenereteSpecificCfgSetPrivateAndPublicThatMatch(t *testing.T) {
	genericCfg := NewMockSignerConfig(privateKey, publicKey)
	_, err := NewMockConfig(genericCfg)
	require.NoError(t, err)
}

func TestGenereteSpecificCfgSetPrivateAndPublicThatMissmatch(t *testing.T) {
	genericCfg := NewMockSignerConfig(privateKey, publicKey2)
	_, err := NewMockConfig(genericCfg)
	require.Error(t, err)
}
