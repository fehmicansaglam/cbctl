package get

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/cbctl/cmd/utils"
	"github.com/fehmicansaglam/cbctl/couchbase"
	"github.com/fehmicansaglam/cbctl/output"
	"github.com/spf13/cobra"
)

var getBucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "Get Couchbase buckets",
	Example: utils.TrimAndIndent(`
	# Retrieve all buckets.
	cbctl get buckets
	`),
	Run: func(cmd *cobra.Command, args []string) {
		handleBucketsLogic()
	},
}

func init() {
	// getBucketsCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
}

var bucketColumns = []output.ColumnDef{
	{Header: "NAME", Type: output.Text},
	{Header: "TYPE", Type: output.Text},
	{Header: "ITEM COUNT", Type: output.Number},
	{Header: "RAM QUOTA (MB)", Type: output.Number},
	{Header: "REPLICA COUNT", Type: output.Number},
	{Header: "FLUSH ENABLED", Type: output.Text},
}

func handleBucketsLogic() {
	buckets, err := couchbase.GetBuckets()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve buckets:", err)
		os.Exit(1)
	}

	columnDefs := bucketColumns
	data := [][]string{}

	for _, bucket := range buckets {
		rowData := map[string]string{
			"NAME":           bucket.Name,
			"TYPE":           bucket.BucketType,
			"RAM QUOTA (MB)": fmt.Sprintf("%d", bucket.RAMQuotaMB),
			"REPLICA COUNT":  fmt.Sprintf("%d", bucket.NumReplicas),
			"FLUSH ENABLED":  fmt.Sprintf("%t", bucket.FlushEnabled),
			"ITEM COUNT":     fmt.Sprintf("%d", bucket.BasicStats.ItemCount),
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	// Sort if specified
	if len(flagSortBy) > 0 {
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "NAME")
	}
}
