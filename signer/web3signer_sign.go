package signer

import (
	"context"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	web3signerclient "github.com/agglayer/go_signer/signer/web3signer_client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	FieldAddress = "address"
	FieldURL     = "url"
)

type Web3SignerClienter interface {
	EthAccounts(ctx context.Context) ([]common.Address, error)
	SignHash(ctx context.Context, address common.Address, hashToSign common.Hash) ([]byte, error)
	SignTx(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error)
}

type Web3SignerConfig struct {
	// URL is the url of the web3 signer
	URL string
	// Address is the address of the account to use, if not specified the first account (if only 1 exposed) will be used
	Address common.Address
}

func NewWeb3SignerConfig(cfg SignerConfig) (Web3SignerConfig, error) {
	var addr common.Address
	addrField, ok := cfg.Config[FieldAddress]
	// Field Address is optional
	if ok {
		s, ok := addrField.(string)
		if !ok {
			return Web3SignerConfig{}, fmt.Errorf("config %s: field %s is not string %v",
				MethodWeb3Signer, FieldAddress, addrField)
		}
		if !common.IsHexAddress(s) {
			return Web3SignerConfig{}, fmt.Errorf("config %s: invalid field %s: %s", MethodWeb3Signer, FieldAddress, s)
		}
		addr = common.HexToAddress(s)
	}
	urlIntf, ok := cfg.Config[FieldURL]
	// Field URL is mandatory
	if !ok {
		return Web3SignerConfig{}, fmt.Errorf("config %s: field %s is not present", MethodWeb3Signer, FieldURL)
	}
	urlStr, ok := urlIntf.(string)
	if !ok {
		return Web3SignerConfig{}, fmt.Errorf("config %s: field %s is not string %v",
			MethodWeb3Signer, FieldURL, cfg.Config["url"])
	}
	return Web3SignerConfig{
		URL:     urlStr,
		Address: addr,
	}, nil
}

type Web3SignerSign struct {
	name    string
	logger  signercommon.Logger
	client  Web3SignerClienter
	address common.Address
}

func NewWeb3SignerSign(name string, logger signercommon.Logger, client Web3SignerClienter,
	address common.Address) *Web3SignerSign {
	return &Web3SignerSign{
		name:    name,
		logger:  logger,
		client:  client,
		address: address,
	}
}

func NewWeb3SignerSignFromConfig(name string, logger signercommon.Logger, cfg Web3SignerConfig) *Web3SignerSign {
	client := web3signerclient.NewWeb3SignerClient(cfg.URL)
	return NewWeb3SignerSign(name, logger, client, cfg.Address)
}

func (e *Web3SignerSign) Initialize(ctx context.Context) error {
	if e.client == nil {
		return fmt.Errorf("%s client is nil", e.logPrefix())
	}
	if e.logger == nil {
		return fmt.Errorf("%s logger is nil", e.logPrefix())
	}
	var zeroAddr common.Address
	if e.address == zeroAddr {
		accounts, err := e.client.EthAccounts(ctx)
		if err != nil {
			return fmt.Errorf("%s error getting ethAccounts to define default public Address. Err:%w",
				e.logPrefix(), err)
		}
		if len(accounts) == 0 {
			return fmt.Errorf("%s no accounts found", e.logPrefix())
		}
		if len(accounts) > 1 {
			return fmt.Errorf("%s more than one account found, please specify the account", e.logPrefix())
		}
		e.logger.Infof("%s Using account %v", e.logPrefix(), accounts[0])
		e.address = accounts[0]
	}
	// Verify it
	_, err := e.SignHash(ctx, common.Hash{})
	if err != nil {
		return fmt.Errorf("%s error verifying signer: sign hash fails: %w", e.logPrefix(), err)
	}
	return nil
}

func (e *Web3SignerSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	var zeroAddr common.Address
	if e.address == zeroAddr {
		return nil, fmt.Errorf("%s no Publicaddress set. Call Initialize first", e.logPrefix())
	}

	return e.client.SignHash(ctx, e.address, hash)
}

func (e *Web3SignerSign) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	return e.client.SignTx(ctx, e.address, tx)
	//return nil, fmt.Errorf("not implemented")
}

func (e *Web3SignerSign) PublicAddress() common.Address {
	return e.address
}

func (e *Web3SignerSign) logPrefix() string {
	return fmt.Sprintf("signer: %s[%s]: ", MethodWeb3Signer, e.name)
}

func (e *Web3SignerSign) String() string {
	return fmt.Sprintf("signer: %s[%s]: pubAddr: %s", MethodWeb3Signer, e.name, e.address.String())
}
