package project

import (
	"fmt"
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
		fmt.Println(settings)
		c.SetOptions(*p, settings)

	},
}

func init() {
	settingsCmd.Flags().BoolVarP(&flagRemoveSourceBranch, "remove-source-branch", "s", false, "update remove source branch value")
	projectCmd.AddCommand(settingsCmd)
}
