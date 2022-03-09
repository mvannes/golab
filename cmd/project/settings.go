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
var flagCommitEventsWillUpdateJira bool
var flagJiraUserPassword string

var jiraSettingsCmd = &cobra.Command{
	Use:   "jira-settings",
	Short: "set jira settings for project. Will ALWAYS change your password",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		p, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}

		settings := gitlab.ProjectJiraSettings{
			Password: flagJiraUserPassword,
		}
		if cmd.Flag("commits-update-jira").Changed {
			settings.CommitEventsUpdateJira = &flagCommitEventsWillUpdateJira
		}

		err = c.UpdateJiraIntegration(*p, settings)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

var jiraSettingsForNamespaceCmd = &cobra.Command{
	Use:   "jira-settings-for-namespace",
	Short: "set settings for all projects in a namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		projects, err := c.Projects(args[0])
		if nil != err {
			log.Fatal(err.Error())
		}

		settings := gitlab.ProjectJiraSettings{
			Password: flagJiraUserPassword,
		}
		if cmd.Flag("commits-update-jira").Changed {
			settings.CommitEventsUpdateJira = &flagCommitEventsWillUpdateJira
		}

		for _, p := range projects {
			if p.Archived {
				continue
			}

			err = c.UpdateJiraIntegration(p, settings)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	},
}

func init() {
	settingsCmd.Flags().BoolVarP(&flagRemoveSourceBranch, "remove-source-branch", "s", false, "update remove source branch value")
	settingsForNamespaceCmd.Flags().BoolVarP(&flagRemoveSourceBranch, "remove-source-branch", "s", false, "update remove source branch value")

	jiraSettingsCmd.Flags().BoolVarP(&flagCommitEventsWillUpdateJira, "commits-update-jira", "c", false, "update commit events update jira value")
	jiraSettingsCmd.Flags().StringVarP(&flagJiraUserPassword, "jira-password", "p", "", "MUST PROVIDE, the jira user password to update to.")
	jiraSettingsCmd.MarkFlagRequired("jira-password")

	jiraSettingsForNamespaceCmd.Flags().BoolVarP(&flagCommitEventsWillUpdateJira, "commits-update-jira", "c", false, "update commit events update jira value")
	jiraSettingsForNamespaceCmd.Flags().StringVarP(&flagJiraUserPassword, "jira-password", "p", "", "MUST PROVIDE, the jira user password to update to.")
	jiraSettingsForNamespaceCmd.MarkFlagRequired("jira-password")

	projectCmd.AddCommand(settingsCmd)
	projectCmd.AddCommand(settingsForNamespaceCmd)
	projectCmd.AddCommand(jiraSettingsCmd)
	projectCmd.AddCommand(jiraSettingsForNamespaceCmd)
}
