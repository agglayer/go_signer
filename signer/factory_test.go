package signer

import (
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/stretchr/testify/require"
)

func TestNewSigner(t *testing.T) {
	logger := log.WithFields("test", "test")
	ctx := context.TODO()

	tests := []struct {
		name             string
		config           signertypes.SignerConfig
		expectedError    bool
		errorMsgContains string
	}{
		{
			name:          "unknown signer method",
			config:        signertypes.SignerConfig{Method: "unknown_method"},
			expectedError: true,
		},
		{
			name:          "empty method is local",
			config:        signertypes.SignerConfig{},
			expectedError: false,
		},
		{
			name: "wrong local config",
			config: signertypes.SignerConfig{
				Config: map[string]interface{}{
					FieldPath: 1234,
				},
			},
			expectedError: true,
		},
		{
			name: "wrong remote config",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodRemoteSigner,
				Config: map[string]interface{}{
					FieldAddress: 1234,
				},
			},
			expectedError:    true,
			errorMsgContains: FieldAddress,
		},
		{
			name: "wrong remote config2",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodRemoteSigner,
				Config: map[string]interface{}{
					FieldAddress: "NOTHEXA",
				},
			},
			expectedError:    true,
			errorMsgContains: FieldAddress,
		},
		{
			name: "wrong remote config3",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodRemoteSigner,
				Config: map[string]interface{}{
					FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
					FieldURL:     1234,
				},
			},
			expectedError:    true,
			errorMsgContains: FieldURL,
		},
		{
			name: "wrong remote missing URL",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodRemoteSigner,
				Config: map[string]interface{}{
					FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
				},
			},
			expectedError:    true,
			errorMsgContains: FieldURL,
		},
		{
			name: "remote config",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodRemoteSigner,
				Config: map[string]interface{}{
					FieldAddress: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
					FieldURL:     "http://localhost:9001",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut, err := NewSigner(ctx, 1, tt.config, "test", logger)
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorMsgContains != "" {
					require.Contains(t, err.Error(), tt.errorMsgContains)
				}
				require.Nil(t, sut)
			} else {
				require.NoError(t, err)
				require.NotNil(t, sut)
			}
		})
	}
}
