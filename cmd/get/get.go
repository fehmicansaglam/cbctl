package get

import (
	"github.com/fehmicansaglam/cbctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Couchbase entities",
	Long: utils.Trim(`
The 'get' command allows you to retrieve information about Couchbase entities.

Available Entities:
  - buckets: List all nodes in the Couchbase cluster.`),
	Example: utils.TrimAndIndent(`
#Retrieve a list of all nodes in the Couchbase cluster.
cbctl get buckets`),
}

func init() {
	getCmd.PersistentFlags().StringSliceVarP(&flagSortBy, "sort-by", "s", []string{}, "Columns to sort by (comma-separated)")
	getCmd.PersistentFlags().StringSliceVarP(&flagColumns, "columns", "c", []string{}, "Columns to display (comma-separated) or 'all'")

	getCmd.AddCommand(getBucketsCmd)
}

func Cmd() *cobra.Command {
	return getCmd
}
