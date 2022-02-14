package project

import "github.com/spf13/cobra"

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project related queries",
}

func Add(root *cobra.Command) {
	root.AddCommand(projectCmd)
}
