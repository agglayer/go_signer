package types

import (
	"bytes"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const (
	configLocal = `
	Signer = {Method = "local", Path = "/app/sequencer.keystore", Password = "test"}
	`
	configRemoteSigner = `
	Signer = {Method = "remote", URL = "http://localhost:8545", Address = "0x1234567890abcdef"}
	`
	configMockPrivateKye = `
	{
     Method: "mock",
     PrivateKey: "0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0",
}
	 `
	configEmpty = `
	Signer = {}
	`
)

func TestUnmarshalLocalConfig(t *testing.T) {
	cfg := struct {
		Signer SignerConfig `jsonschema:"omitempty" mapstructure:"Signer"`
	}{}
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configLocal)))
	require.NoError(t, err)
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}
	err = viper.Unmarshal(&cfg, decodeHooks...)
	require.NoError(t, err)
	require.Equal(t, MethodLocal, cfg.Signer.Method)
	require.Equal(t, "/app/sequencer.keystore", cfg.Signer.Config["path"])
	require.Equal(t, "test", cfg.Signer.Config["password"])
}

func TestUnmarshalRemoteSignerConfig(t *testing.T) {
	cfg := struct {
		Signer SignerConfig `jsonschema:"omitempty" mapstructure:"Signer"`
	}{}
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configRemoteSigner)))
	require.NoError(t, err)
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}
	err = viper.Unmarshal(&cfg, decodeHooks...)
	require.NoError(t, err)
	require.Equal(t, MethodRemoteSigner, cfg.Signer.Method)
	require.Equal(t, "http://localhost:8545", cfg.Signer.Config["url"])
	require.Equal(t, "0x1234567890abcdef", cfg.Signer.Config["address"])
}

func TestUnmarshalConfig(t *testing.T) {
	cfg := struct {
		Signer SignerConfig `jsonschema:"omitempty" mapstructure:"Signer"`
	}{}
	viper.SetConfigType("toml")
	{

	}
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configRemoteSigner)))
	require.NoError(t, err)
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}
	err = viper.Unmarshal(&cfg, decodeHooks...)
	require.NoError(t, err)
	require.Equal(t, MethodRemoteSigner, cfg.Signer.Method)
	require.Equal(t, "http://localhost:8545", cfg.Signer.Config["url"])
	require.Equal(t, "0x1234567890abcdef", cfg.Signer.Config["address"])
}

func TestUnmarshalEmptyConfig(t *testing.T) {
	cfg := struct {
		Signer SignerConfig `jsonschema:"omitempty" mapstructure:"Signer"`
	}{}
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configEmpty)))
	require.NoError(t, err)
	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}
	err = viper.Unmarshal(&cfg, decodeHooks...)
	require.NoError(t, err)
	require.Equal(t, "", cfg.Signer.Method.String())
	require.Equal(t, 0, len(cfg.Signer.Config))
}

func TestSignerConfigString(t *testing.T) {
	cfg := SignerConfig{
		Method: MethodLocal,
		Config: map[string]any{
			"Path":     "/app/sequencer.keystore",
			"Password": "test",
		},
	}
	expected := "SignerConfig:Method: local\n Config[Password]: test\n Config[Path]: /app/sequencer.keystore\n"
	require.Equal(t, expected, cfg.String())

	cfg = SignerConfig{
		Method: MethodRemoteSigner,
		Config: map[string]any{
			"URL":     "http://localhost:8545",
			"Address": "0x1234567890abcdef",
		},
	}
	expected = "SignerConfig:Method: remote\n Config[Address]: 0x1234567890abcdef\n Config[URL]: http://localhost:8545\n"
	require.Equal(t, expected, cfg.String())

	require.Equal(t, "SignerConfig:Method: \n", SignerConfig{}.String())
	require.Equal(t, "SignerConfig:Method: test\n Config[key]: value\n Config[key2]: 4\n", SignerConfig{
		Method: "test",
		Config: map[string]any{"key": "value",
			"key2": 4},
	}.String())
	require.Equal(t, "SignerConfig:Method: test-nil\n", SignerConfig{
		Method: "test-nil",
		Config: nil,
	}.String())
}
