package signer

import (
	"fmt"
	"testing"

	"github.com/agglayer/go_signer/log"
	"github.com/stretchr/testify/require"
)

func TestMockSignRandomPrivateKey(t *testing.T) {
	sut, err := NewMockSign("test-mock-signer",
		log.WithFields("module", "test"),
		NewMockSignerConfig("", ""),
		1774)
	require.NoError(t, err, "should not return error when creating mock signer with random private key")
	fmt.Print("public addr: ", sut.PublicAddress())
}

func TestMockSignRandomPrivateKey2(t *testing.T) {
}
