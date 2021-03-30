package cmd

import (
	"fmt"
	"github.com/ariarijp/retag/internal/redash"
	"github.com/go-resty/resty/v2"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"io/ioutil"
	"regexp"
	"strings"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff of YAML file on your local storage and queries on the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		yamlFilePath, _ := cmd.Flags().GetString("yaml")

		client := resty.New()
		redash.Init(baseUrl, apiKey, client)

		local, err := ioutil.ReadFile(yamlFilePath)
		if err != nil {
			return xerrors.Errorf("IO error: %+w", err)
		}

		remote, err := redash.ExportQueriesAsYAML()
		if err != nil {
			return err
		}

		fmt.Println(diff(remote.Bytes(), local))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().String("yaml", "", "YAML")
	diffCmd.MarkFlagRequired("yaml")
}

func diff(s1, s2 []byte) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(string(s1), string(s2))
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)

	return DiffPlainText(diffs)
}

func DiffPlainText(diffs []diffmatchpatch.Diff) string {
	r, _ := regexp.Compile(`(?m)^`)
	s := ""
	for _, d := range diffs {
		tmp := strings.TrimRight(d.Text, "\n")
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			s = s + r.ReplaceAllString(tmp, "-") + "\n"
		case diffmatchpatch.DiffEqual:
			s = s + r.ReplaceAllString(tmp, " ") + "\n"
		case diffmatchpatch.DiffInsert:
			s = s + r.ReplaceAllString(tmp, "+") + "\n"
		}
	}

	return s
}