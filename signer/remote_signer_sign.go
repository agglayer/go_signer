package signer

import (
	"context"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	web3signerclient "github.com/agglayer/go_signer/signer/remotesignerclient"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	FieldAddress = "address"
	FieldURL     = "url"
)

var zeroAddr common.Address

type RemoteSignerClienter interface {
	EthAccounts(ctx context.Context) ([]common.Address, error)
	SignHash(ctx context.Context, address common.Address, hashToSign common.Hash) ([]byte, error)
	SignTx(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error)
}

type RemoteSignerConfig struct {
	// URL is the url of the web3 signer
	URL string
	// Address is the address of the account to use, if not specified the first account (if only 1 exposed) will be used
	Address common.Address
}

func NewRemoteSignerConfig(cfg signertypes.SignerConfig) (RemoteSignerConfig, error) {
	var addr common.Address
	addrField, ok := cfg.Config[FieldAddress]
	// Field Address is optional
	if ok {
		s, ok := addrField.(string)
		if !ok {
			return RemoteSignerConfig{}, fmt.Errorf("config %s: field %s is not string %v",
				signertypes.MethodRemoteSigner, FieldAddress, addrField)
		}
		if !common.IsHexAddress(s) {
			return RemoteSignerConfig{},
				fmt.Errorf("config %s: invalid field %s: %s", signertypes.MethodRemoteSigner, FieldAddress, s)
		}
		addr = common.HexToAddress(s)
	}
	urlIntf, ok := cfg.Config[FieldURL]
	// Field URL is mandatory
	if !ok {
		return RemoteSignerConfig{},
			fmt.Errorf("config %s: field %s is not present", signertypes.MethodRemoteSigner, FieldURL)
	}
	urlStr, ok := urlIntf.(string)
	if !ok {
		return RemoteSignerConfig{}, fmt.Errorf("config %s: field %s is not string %v",
			signertypes.MethodRemoteSigner, FieldURL, cfg.Config["url"])
	}
	return RemoteSignerConfig{
		URL:     urlStr,
		Address: addr,
	}, nil
}

type RemoteSignerSign struct {
	name    string
	logger  signercommon.Logger
	client  RemoteSignerClienter
	address common.Address
}

func NewRemoteSignerSign(name string, logger signercommon.Logger, client RemoteSignerClienter,
	address common.Address) *RemoteSignerSign {
	return &RemoteSignerSign{
		name:    name,
		logger:  logger,
		client:  client,
		address: address,
	}
}

func NewRemoteSignerSignFromConfig(name string, logger signercommon.Logger, cfg RemoteSignerConfig) *RemoteSignerSign {
	client := web3signerclient.NewRemoteSignerClient(cfg.URL)
	return NewRemoteSignerSign(name, logger, client, cfg.Address)
}

func (e *RemoteSignerSign) Initialize(ctx context.Context) error {
	if e.client == nil {
		return fmt.Errorf("%s client is nil", e.logPrefix())
	}
	if e.logger == nil {
		return fmt.Errorf("%s logger is nil", e.logPrefix())
	}
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
	return nil
}

func (e *RemoteSignerSign) SignHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("remote eth_sign use EIP155 that change the hash to sign. So you can't use to sign")
}

func (e *RemoteSignerSign) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	return e.client.SignTx(ctx, e.address, tx)
}

func (e *RemoteSignerSign) PublicAddress() common.Address {
	return e.address
}

func (e *RemoteSignerSign) logPrefix() string {
	return fmt.Sprintf("signer: %s[%s]: ", signertypes.MethodRemoteSigner, e.name)
}

func (e *RemoteSignerSign) String() string {
	return fmt.Sprintf("signer: %s[%s]: pubAddr: %s", signertypes.MethodRemoteSigner, e.name, e.address.String())
}
