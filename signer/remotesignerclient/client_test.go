package remotesignerclient

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestWeb3SignerClientUsingDockerNoKey(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	sut := NewRemoteSignerClient("http://localhost:9999")
	require.NotNil(t, sut)
	_, err := sut.SignHash(context.Background(), common.Address{}, common.Hash{})
	require.Error(t, err, "sign address doesn't have a private key to sign")
	addrs, err := sut.EthAccounts(context.Background())
	require.NoError(t, err)
	require.True(t, len(addrs) > 0)
}
