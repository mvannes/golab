package project

import (
	"fmt"
	"log"

	"github.com/mvannes/golab/gitlab"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in a namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		c := gitlab.NewClient()
		projects, err := c.Projects(args[0])
		if nil != err {
			log.Fatal(err.Error())
		}
		for _, p := range projects {
			fmt.Println(p.NameWithNamespace)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get single project by namespace and name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		c := gitlab.NewClient()
		p, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}

		fmt.Println(p.NameWithNamespace)
	},
}

func init() {
	projectCmd.AddCommand(getCmd)
}
