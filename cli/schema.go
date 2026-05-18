package cli

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gitlayzer/seapack/core/config"
	"github.com/urfave/cli/v3"
)

var SchemaCommand = &cli.Command{
	Name:                  "schema",
	Usage:                 "outputs the JSON schema for the SeaPack config",
	EnableShellCompletion: true,
	Flags:                 []cli.Flag{},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		schema := config.GetJsonSchema()

		schemaJson, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return cli.Exit(err, 1)
		}

		_, _ = os.Stdout.Write(schemaJson)
		_, _ = os.Stdout.Write([]byte("\n"))

		return nil
	},
}
