package version

import (
	"os"

	"github.com/agglayer/go_signer"
	cli "github.com/urfave/cli/v2"
)

func VersionCmd(*cli.Context) error {
	go_signer.PrintVersion(os.Stdout)
	return nil
}
