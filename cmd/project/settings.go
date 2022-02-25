package project

import (
	"log"

	"github.com/mvannes/golab/gitlab"
	"github.com/spf13/cobra"
)

var flagRemoveSourceBranch bool

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "set settings for project",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		p, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		settings := gitlab.ProjectSettings{}
		if cmd.Flag("remove-source-branch").Changed {
			settings.RemoveSourceBranchAfterMerge = &flagRemoveSourceBranch
		}
		c.SetOptions(*p, settings)
	},
}

var settingsForNamespaceCmd = &cobra.Command{
	Use:   "settings-for-namespace",
	Short: "set settings for all projects in a namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		projects, err := c.Projects(args[0])
		if nil != err {
			log.Fatal(err.Error())
		}

		for _, p := range projects {
			settings := gitlab.ProjectSettings{}
			if cmd.Flag("remove-source-branch").Changed && flagRemoveSourceBranch != p.RemoveSourceBranchAfterMerge {
				settings.RemoveSourceBranchAfterMerge = &flagRemoveSourceBranch
			}
			c.SetOptions(p, settings)
		}
	},
}

func init() {
	settingsCmd.Flags().BoolVarP(&flagRemoveSourceBranch, "remove-source-branch", "s", false, "update remove source branch value")
	settingsForNamespaceCmd.Flags().BoolVarP(&flagRemoveSourceBranch, "remove-source-branch", "s", false, "update remove source branch value")
	projectCmd.AddCommand(settingsCmd)
	projectCmd.AddCommand(settingsForNamespaceCmd)
}
