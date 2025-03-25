package signer

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
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
	mockRemoteSignerClient.EXPECT().EthAccounts(ctx).Return([]common.Address{
		common.HexToAddress("0x1234"),
	}, nil)
	mockRemoteSignerClient.EXPECT().SignHash(ctx, common.HexToAddress("0x1234"), mock.Anything).Return([]byte{}, nil).Once()
	err := sut.Initialize(ctx)
	require.NoError(t, err)
	signData := []byte{0x01, 0x02, 0x03}
	mockRemoteSignerClient.EXPECT().SignHash(ctx, mock.Anything, mock.Anything).Return(signData, nil)
	sign, err := sut.SignHash(ctx, common.HexToHash("0x1234"))
	require.NoError(t, err)
	require.NotNil(t, sign)
	require.Equal(t, signData, sign)
}
