package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// shutdownDockerAfterTest:  set to false if you want to inspect the container
// after running the test
const shutdownDockerAfterTest = true

// dockerIsAlreadyRunning: set to true if you want to start manually the containers
// or you want to take advantage of previous run
const dockerIsAlreadyRunning = false

const gethURL = "http://localhost:8545"
const defaultChainID = uint64(1337)

func testSendEthTx(t *testing.T, fromAddress common.Address, txSigner signer.TxSigner) {
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
