package e2e

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	gethURL           = "http://localhost:8545"
	defaultChainID    = uint64(1337)
	publicAddressTest = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
)

type e2eTestParams struct {
	createSignerFunc      func(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error)
	canSign               bool
	expectedPublicAddress string
}

func createLocalSigner(t *testing.T, ctx context.Context, chainID uint64) (signertypes.Signer, error) {
	t.Helper()
	password, err := os.ReadFile("key_store/funded_addr.password")
	require.NoError(t, err)
	trimmedPassword := strings.TrimSpace(string(password))
	return signer.NewSigner(ctx, chainID, signertypes.SignerConfig{
		Method: signertypes.MethodLocal,
		Config: map[string]interface{}{
			signer.FieldPath:     "key_store/funded_addr",
			signer.FieldPassword: trimmedPassword,
		},
	}, "test", log.WithFields("module", "test"))
}

func checkSignatureAgainstReferenceLocalSigner(t *testing.T,
	ctx context.Context,
	hashToSign common.Hash, signature []byte) {
	t.Helper()
	localSigner, err := createLocalSigner(t, ctx, 0)
	require.NoError(t, err)
	localSign, ok := localSigner.(*signer.LocalSign)
	require.True(t, ok)
	err = localSigner.Initialize(ctx)
	require.NoError(t, err)
	err = localSign.Verify(hashToSign, signature[0:64])
	require.NoError(t, err)
	t.Log("signature: ", common.Bytes2Hex(signature))
	t.Log("signature length: ", len(signature))
	t.Log("signature without V: ", common.Bytes2Hex(signature[0:64]))
	t.Log("signature V: ", signature[64])
	t.Log("✔️ signature Verified!")
}

func testGenericSignerE2E(t *testing.T, params e2eTestParams) {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}

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
	if params.expectedPublicAddress != "" {
		require.Equal(t, params.expectedPublicAddress, sign.PublicAddress().String())
	} else {
		require.Equal(t, publicAddressTest, sign.PublicAddress().String())
	}
	require.NoError(t, err)
	hashToSign := common.Hash{}
	signed, err := sign.SignHash(ctx, hashToSign)
	if params.canSign {
		require.NoError(t, err)
		require.NotNil(t, signed)
		log.Debugf("signed hash: %s", common.Bytes2Hex(signed))
		checkSignatureAgainstReferenceLocalSigner(t, ctx, hashToSign, signed)
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
