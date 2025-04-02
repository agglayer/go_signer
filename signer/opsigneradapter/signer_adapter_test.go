package opsigneradapter

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestConvertPublicKeyToAddress(t *testing.T) {
	publicKey := common.Hex2Bytes("048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5")
	expectedAddress := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	address := convertPublicKeyToAddress(publicKey)
	require.Equal(t, expectedAddress, strings.ToLower(address.String()))
}
