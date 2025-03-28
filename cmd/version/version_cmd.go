package version

import (
	"os"

	gosigner "github.com/agglayer/go_signer"
	cli "github.com/urfave/cli/v2"
)

func VersionCmd(*cli.Context) error {
	gosigner.PrintVersion(os.Stdout)
	return nil
}
