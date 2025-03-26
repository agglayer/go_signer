package remotesignerclient

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/test/e2e/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestWeb3SignerClientUsingDockerNoKey(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	helpers.StartWeb3SignerDocker(t)
	defer func() {
		helpers.StopWeb3SignerDocker(t)
	}()
	sut := NewRemoteSignerClient("http://localhost:9999")
	require.NotNil(t, sut)
	_, err := sut.SignHash(context.Background(), common.Address{}, common.Hash{})
	require.Error(t, err)
	addrs, err := sut.EthAccounts(context.Background())
	require.Error(t, err)
	require.Nil(t, addrs)
}
