package remotesignerclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/agglayer/go_signer/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type RemoteSignerClient struct {
	url string
}

func NewRemoteSignerClient(url string) *RemoteSignerClient {
	return &RemoteSignerClient{
		url: url,
	}
}

func (e *RemoteSignerClient) EthAccounts(ctx context.Context) ([]common.Address, error) {
	response, err := rpc.JSONRPCCallWithContext(ctx, e.url, "eth_accounts")
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%v %v", response.Error.Code, response.Error.Message)
	}
	var result []common.Address
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *RemoteSignerClient) SignHash(ctx context.Context,
	address common.Address,
	hashToSign common.Hash) ([]byte, error) {
	params := []interface{}{address, hashToSign}
	response, err := rpc.JSONRPCCallWithContext(ctx, e.url, "eth_sign", params...)
	if err != nil {
		return nil, fmt.Errorf("signHash eth_sign RPC call fails. Err: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("signHash fails. Code:%v Message:%v", response.Error.Code, response.Error.Message)
	}
	var resultStr string
	err = json.Unmarshal(response.Result, &resultStr)
	if err != nil {
		return nil, err
	}
	result := common.FromHex(resultStr)
	return result, nil
}

func (e *RemoteSignerClient) SignTx(ctx context.Context,
	from common.Address, tx *types.Transaction) (*types.Transaction, error) {
	params := map[string]interface{}{
		"from": from.String(),
	}
	if tx.To() != nil {
		params["to"] = tx.To().String()
	}
	if tx.Gas() != 0 {
		params["gas"] = tx.Gas()
	}
	if tx.GasPrice() != nil {
		params["gasPrice"] = tx.GasPrice().String()
	}
	params["nonce"] = tx.Nonce()

	if tx.Value() != nil {
		params["value"] = tx.Value().String()
	}
	if tx.Data() != nil {
		params["data"] = tx.Data()
	}
	// Fields maxPriorityFeePerGas and maxFeePerGas are not set because the API doesn't support them:
	// - https://docs.web3signer.consensys.io/reference/api/json-rpc#eth_signtransaction
	// - https://www.quicknode.com/docs/ethereum/eth_signTransaction
	response, err := rpc.JSONRPCCallWithContext(ctx, e.url, "eth_signTransaction", params)
	if err != nil {
		return nil, fmt.Errorf("SignTx eth_signTransaction RPC call fails. Err: %w", err)
	}
	if response.Error != nil {
		return nil, fmt.Errorf("SignTx fails. Code:%v Message:%v", response.Error.Code, response.Error.Message)
	}
	var resultStr string
	if err = json.Unmarshal(response.Result, &resultStr); err != nil {
		return nil, fmt.Errorf("SignTx unmarshal fails. Err: %w", err)
	}
	log.Debugf("SignTx result: (%d) %s", len(resultStr), resultStr)
	encodedTx := common.FromHex(resultStr)
	var resTx *types.Transaction
	if err = rlp.DecodeBytes(encodedTx, &resTx); err != nil {
		return nil, fmt.Errorf("SignTx rlp.DecodeBytes fails. Err: %w", err)
	}
	signer := types.NewEIP155Signer(resTx.ChainId())
	// sanity check:  Just verify the signingHash
	if signer.Hash(resTx) != signer.Hash(tx) {
		return nil, fmt.Errorf("SignTx signingHash differs:  %s!=%s", signer.Hash(tx).String(), signer.Hash(resTx).String())
	}
	return resTx, err
}
