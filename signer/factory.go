package signer

import (
	"context"
	"fmt"

	signercommon "github.com/agglayer/go_signer/common"
	"github.com/agglayer/go_signer/signer/opsigneradapter"
	"github.com/agglayer/go_signer/signer/types"
)

var (
	ErrUnknownSignerMethod = fmt.Errorf("unknown signer method")
)

func NewSigner(ctx context.Context, chainID uint64, cfg types.SignerConfig, name string,
	logger signercommon.Logger) (types.Signer, error) {
	var res types.Signer
	var err error
	if cfg.Method == "" {
		logger.Warnf("No signer method specified, defaulting to local (keystore file)")
		cfg.Method = types.MethodLocal
	}
	switch cfg.Method {
	case types.MethodNone:
		res = &NoneSign{}
	case types.MethodLocal:
		specificCfg, err := NewLocalConfig(cfg)
		if err != nil {
			return nil, err
		}
		res = NewLocalSign(name, logger, specificCfg, chainID)
	case types.MethodRemoteSigner:
		specificCfg, err := NewRemoteSignerConfig(cfg)
		if err != nil {
			return nil, err
		}
		res = NewRemoteSignerSignFromConfig(name, logger, specificCfg)
	case types.MethodGCPKMS:
		res, err = opsigneradapter.NewSignerAdapterFromConfig(ctx, logger, cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown signer method %s", cfg.Method)
	}
	return res, nil
}
