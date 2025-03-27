package e2e

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/agglayer/go_signer/log"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/agglayer/go_signer/test/e2e/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (

	// dockerIsAlreadyRunning: set to true if you want to start manually the containers
	// or you want to take advantage of previous run
	dockerIsAlreadyRunning = false

	gethURL                = "http://localhost:8545"
	defaultChainID         = uint64(1337)
	publicAddressTest      = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	sleepWaitDockerHealthy = 40 * time.Second
)

var (
	mu         sync.Mutex
	dockerIsUp = false
)

type e2eTestParams struct {
	createSignerFunc func(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error)
	canSign          bool
}

func testGenericSignerE2E(t *testing.T, params e2eTestParams) {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}
	mu.Lock()
	if !dockerIsUp && !dockerIsAlreadyRunning {
		dockerCompose := helpers.NewDockerCompose()
		dockerCompose.Down(t)
		dockerCompose.Up(t)
		dockerCompose.WaitHealthy(t, sleepWaitDockerHealthy)
		dockerIsUp = true
	}
	mu.Unlock()
	ctx := context.TODO()
	ethClient, err := ethclient.Dial(gethURL)
	require.NoError(t, err)
	defer ethClient.Close()
	chainID, err := ethClient.ChainID(ctx)
	require.NoError(t, err)
	log.Info("chainID: ", chainID.Uint64())
	sign, err := params.createSignerFunc(t, ctx, chainID.Uint64())
	require.NoError(t, err)

	err = sign.Initialize(ctx)
	require.NoError(t, err)
	require.Equal(t, publicAddressTest, sign.PublicAddress().String())

	signed, err := sign.SignHash(ctx, common.Hash{})
	if params.canSign {
		require.NoError(t, err)
		require.NotNil(t, signed)
		log.Debugf("signed hash: %s", common.Bytes2Hex(signed))
		require.Equal(t, "b8823364c90ea0d2700d5ad0fe39d16778bc07ce7df4779ff35e4b2660d043cb74a002439225d1d518f9f1cf3db005f5e143196543fd5146a34bf63f0b810ade00",
			common.Bytes2Hex(signed))
	} else {
		require.Error(t, err)
	}

	testSendEthTx(t, sign.PublicAddress(), sign)
}

func testSendEthTx(t *testing.T, fromAddress common.Address, txSigner signertypes.TxSigner) {
	t.Helper()
	client, err := ethclient.Dial(gethURL)
	require.NoError(t, err)
	defer client.Close()
	toAddress := common.HexToAddress("0x1234567890ABCDEF1234567890ABCDEF12345678")
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	require.NoError(t, err)
	// 0.0001 ETH en wei = 0.0001 * 10^18 = 100,000,000,000,000 (1e14)
	value := big.NewInt(1e14) //nolint:mnd
	gasLimit := uint64(21000) //nolint:mnd
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	signedTx, err := txSigner.SignTx(context.Background(), tx)
	require.NoError(t, err)
	err = client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	require.NoError(t, err)
	log.Infof("balance: %s", balance.String())
}
