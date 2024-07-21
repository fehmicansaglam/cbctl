package query

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/cbctl/cmd/utils"
	"github.com/fehmicansaglam/cbctl/couchbase"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query Couchbase",
	Long:  `This command allows you to query Couchbase.`,
	Example: utils.TrimAndIndent(`
cbctl query articles
cbctl query articles --id 61`),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucket := args[0]

		response, err := couchbase.SearchDocuments(bucket, flagId)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to query:", err)
			os.Exit(1)
		}

		fmt.Println(response.Json)
	},
}

func Cmd() *cobra.Command {
	return queryCmd
}

func init() {
	queryCmd.Flags().StringVar(&flagId, "id", "", "Document ID to fetch")
}
