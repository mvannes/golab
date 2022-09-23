package project

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"

	"github.com/mvannes/golab/gitlab"
)

type JobInfo struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	AllowFailure bool   `json:"allow_failure"`
}

var jobsInBranchCmd = &cobra.Command{
	Use:   "jobs-in-branch",
	Short: "List all for the latest pipeline in a project by [namespace], [name] and [branch-or-tag]",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		p, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		pipelines, err := c.Pipelines(*p, args[2])
		if nil != err {
			log.Fatal(err.Error())
		}
		if len(pipelines) == 0 {
			log.Fatal("No pipeline found to resolve jobs for.")
		}
		pipeline := pipelines[0]
		jobs, err := c.JobsForPipeline(*p, *pipeline)
		if err != nil {
			log.Fatal(err)
		}

		result := []JobInfo{}
		for _, j := range jobs {
			result = append(result, JobInfo{
				Name:         j.Name,
				Status:       j.Status,
				AllowFailure: j.AllowFailure,
			})
		}

		b, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	projectCmd.AddCommand(jobsInBranchCmd)
}
