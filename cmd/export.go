package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/ariarijp/retag/internal/redash"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export queries from the server to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		redash.Init(baseUrl, apiKey, client)

		d, err := redash.ExportQueriesAsYAML()
		if err != nil {
			return xerrors.Errorf("Export error: %+w", err)
		}

		fmt.Print(d.String())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
