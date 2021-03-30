package cmd

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/ariarijp/retag/internal/redash"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import queries from YAML file to the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		yamlFilePath, _ := cmd.Flags().GetString("yaml")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		client := resty.New()
		redash.Init(baseUrl, apiKey, client)

		queries, err := loadQueries(yamlFilePath)
		if err != nil {
			return err
		}

		for _, q := range queries {
			time.Sleep(100 * time.Millisecond)

			q2, err := redash.GetQuery(q.Id)
			if err != nil {
				return err
			}

			if reflect.DeepEqual(q, *q2) {
				fmt.Printf("SKIP: Query #%d will not be updated.\n", q.Id)
				continue
			} else if dryRun {
				fmt.Printf("DRYRUN: Query #%d will be updated.\n", q.Id)
				continue
			}

			err = redash.UpdateQuery(q)
			if err != nil {
				return err
			}
			fmt.Printf("RESULT: Query #%d updated.\n", q.Id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String("yaml", "", "YAML")
	importCmd.MarkFlagRequired("yaml")
	importCmd.Flags().Bool("dry-run", true, "Dry run")
}

func loadQueries(p string) ([]redash.Query, error) {
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, xerrors.Errorf("IO error: %+w", err)
	}

	result := map[string][]redash.Query{
		"queries": {},
	}
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
		return nil, xerrors.Errorf("YAML error: %+w", err)
	}

	return result["queries"], nil
}
