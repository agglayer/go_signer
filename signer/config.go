package signer

// SignerConfig is the configuration for the Signer. It's generic because it support
// multiple methods of signing. In order to get yourself familiarized with t
// he exact parameters you must check local.go and remote.go
// Examples of valid values:
// { Method="Local", Path="path/to/keystore", Password="password" }
// { Method="remote", URL="http://localhost:9000", Address="0x1234567890abcdef" }
type SignerConfig struct {
	// Method is the method to use to sign
	Method SignMethod `jsonschema:"enum=local, enum=remote_eth" mapstructure:"Method"`
	// Config is the configuration for the signer (depend on Method field)
	Config map[string]any `jsonschema:"omitempty" mapstructure:",remain"`
}
