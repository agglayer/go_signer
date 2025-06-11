package signer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPrivateKeyHex = "0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0"
	testPublicKeyHex  = "0xc653eCD4AC5153a3700Fb13442Bcf00A691cca16"
)

func TestGenereteSpecificOnlyPrivateKey(t *testing.T) {
	genericCfg := NewMockSignerConfig(testPrivateKeyHex)
	specificCfg, err := NewMockConfig(genericCfg)
	require.NoError(t, err, "should not return error when creating specific config from generic config")
	require.NotNil(t, specificCfg, "specific config should not be nil")
	require.NotNil(t, specificCfg.privateKey, "specific config private key should not be nil")
	require.Equal(t, "MockSignConfigure{mode: PrivateKey, privateKey: SET (public:0xc653eCD4AC5153a3700Fb13442Bcf00A691cca16)}", specificCfg.String())
}

func TestGenereteSpecificCfgNoKeySoRandom(t *testing.T) {
	genericCfg := NewMockSignerConfig("")
	specificCfg, err := NewMockConfig(genericCfg)
	require.NoError(t, err)
	require.Equal(t, "MockSignConfigure{mode: Random}", specificCfg.String())
	require.Nil(t, specificCfg.privateKey)
}
