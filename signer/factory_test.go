package signer

import (
	"bytes"
	"context"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/agglayer/go_signer/signer/opsigneradapter"
	signertypes "github.com/agglayer/go_signer/signer/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
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
			name:          "empty method is local",
			config:        signertypes.SignerConfig{Method: signertypes.MethodNone},
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
		{
			name: "AWS config",
			config: signertypes.SignerConfig{
				Method: signertypes.MethodAWSKMS,
				Config: map[string]interface{}{
					opsigneradapter.FieldKeyName: "path-to-key",
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

func TestNewSignerFromConfigFile(t *testing.T) {
	testcases := []struct {
		name                 string
		configFileContent    []string
		expectedConfigString string
		expectedSignerString string
	}{
		{
			name: "local signer",
			configFileContent: []string{`
			[Signer]
			Method = "local"
			Path = "/path/to/keystore"
			Password = "password"
		`,
				`Signer={ Method="local", Path="/path/to/keystore", Password="password"}`,
			},
			expectedConfigString: "SignerConfig:Method: local\n Config[password]: password\n Config[path]: /path/to/keystore\n",
			expectedSignerString: "signer: local[local signer]:  initialized:false path:{/path/to/keystore password}, pubAddr: ???",
		},
		{
			name: "mock random",
			configFileContent: []string{`
			[Signer]
			Method = "mock"
			`,
				`Signer={ Method="mock"}`,
			},
			expectedConfigString: "SignerConfig:Method: mock\n",
			expectedSignerString: "MockSign{name:mock random, mode:Random, initialized: false}",
		},
		{
			name: "mock privateKey",
			configFileContent: []string{`
			[Signer]
			Method = "mock"
			PrivateKey = "0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0"
			`,
				`Signer={ Method="mock", PrivateKey="0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0"}`,
			},
			expectedConfigString: "SignerConfig:Method: mock\n Config[privatekey]: 0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0\n",
			expectedSignerString: "MockSign{name:mock privateKey, mode:PrivateKey, initialized: false}",
		},
	}

	for _, tc := range testcases {
		for _, content := range tc.configFileContent {
			t.Run(tc.name, func(t *testing.T) {
				cfg := struct {
					Signer signertypes.SignerConfig `jsonschema:"omitempty" mapstructure:"Signer"`
				}{}
				viper.SetConfigType("toml")
				err := viper.ReadConfig(bytes.NewBuffer([]byte(content)))
				require.NoError(t, err)
				decodeHooks := []viper.DecoderConfigOption{
					viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
						mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
				}
				err = viper.Unmarshal(&cfg, decodeHooks...)
				require.NoError(t, err)
				require.Equal(t, tc.expectedConfigString, cfg.Signer.String())
				ctx := context.TODO()
				sut, err := NewSigner(ctx, 1, cfg.Signer, tc.name, log.WithFields("test", "test"))
				require.NoError(t, err)
				require.Equal(t, tc.expectedSignerString, sut.String())
			})
		}
	}
}
