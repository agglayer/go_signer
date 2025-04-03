package signer

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestFailsCantSetAddressToUse(t *testing.T) {
	mockRemoteSignerClient := mocks.NewRemoteSignerClienter(t)
	ctx := context.TODO()
	logger := log.WithFields("test", "test")
	sut := NewRemoteSignerSign("name", logger, mockRemoteSignerClient, common.Address{})
	mockRemoteSignerClient.EXPECT().EthAccounts(ctx).Return([]common.Address{}, nil)
	err := sut.Initialize(ctx)
	require.Error(t, err)
}

func TestInitialize(t *testing.T) {
	sut := RemoteSignerSign{}
	ctx := context.TODO()
	err := sut.Initialize(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "client is nil")
	sut = RemoteSignerSign{
		client: mocks.NewRemoteSignerClienter(t),
	}
	err = sut.Initialize(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "logger is nil")
}

func TestFailsSetAddressToUse(t *testing.T) {
	mockRemoteSignerClient := mocks.NewRemoteSignerClienter(t)
	ctx := context.TODO()
	logger := log.WithFields("test", "test")
	sut := NewRemoteSignerSign("name", logger, mockRemoteSignerClient, common.Address{})
	publicAddr := common.HexToAddress("0x1234")
	mockRemoteSignerClient.EXPECT().EthAccounts(ctx).Return([]common.Address{
		publicAddr,
	}, nil)
	err := sut.Initialize(ctx)
	require.NoError(t, err)
	require.Equal(t, publicAddr, sut.PublicAddress())
}
