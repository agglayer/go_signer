package types

import (
	"fmt"
	"sort"
	"strings"
)

type SignMethod string

var (
	MethodNone         SignMethod = "none"
	MethodLocal        SignMethod = "local"
	MethodRemoteSigner SignMethod = "remote"
	MethodGCPKMS       SignMethod = "GCP"
	MethodAWSKMS       SignMethod = "AWS"
)

func (m SignMethod) String() string {
	return string(m)
}

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

func (c SignerConfig) Get(key string) (string, error) {
	v, ok := c.Config[key]
	if !ok {
		v, ok = c.Config[strings.ToLower(key)]
		if !ok {
			return "", fmt.Errorf("key %s not found", key)
		}
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}
	return s, nil
}

func (c SignerConfig) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Method: %s\n", c.Method))
	keys := make([]string, 0, len(c.Config))
	for k := range c.Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" Config[%s]: %v\n", k, c.Config[k]))
	}
	return sb.String()
}
