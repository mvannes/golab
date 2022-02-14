package mergerequest

import (
	"github.com/spf13/cobra"
)

var mergeRequestCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge request related queries",
}

func Add(root *cobra.Command) {
	root.AddCommand(mergeRequestCmd)
}
