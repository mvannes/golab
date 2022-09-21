package project

import (
	"encoding/json"
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
			if p.Archived == true {
				continue
			}
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

type ProjectVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ProjectVariables struct {
	Project   string            `json:"project"`
	Variables []ProjectVariable `json:"variables"`
}

var variablesCmd = &cobra.Command{
	Use:   "ci-variables-for-namespace",
	Short: "Get project CI variables for all projects in namespace as json",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		c := gitlab.NewClient()
		projects, err := c.Projects(args[0])
		if nil != err {
			log.Fatal(err.Error())
		}
		projectVariables := []ProjectVariables{}

		for _, p := range projects {
			if p.Archived == true {
				continue
			}
			if p.BuildsAccessLevel == "disabled" {
				continue
			}
			variables, err := c.CIVariables(p)
			if nil != err {
				log.Fatal(err.Error())
			}

			ciVariables := []ProjectVariable{}
			for _, v := range variables {
				ciVariables = append(ciVariables, ProjectVariable{Name: v.Key, Value: v.Value})
			}

			projectVariables = append(
				projectVariables,
				ProjectVariables{Project: p.NameWithNamespace, Variables: ciVariables},
			)
		}

		b, err := json.Marshal(projectVariables)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(string(b))
	},
}

var variableCmd = &cobra.Command{
	Use:   "ci-variables",
	Short: "Get single project CI variables by namespace and name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		c := gitlab.NewClient()
		p, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}

		variables, err := c.CIVariables(*p)
		if nil != err {
			log.Fatal(err.Error())
		}
		for _, v := range variables {
			fmt.Println(v.Key, v.Value)
		}
	},
}

func init() {
	projectCmd.AddCommand(getCmd)
	projectCmd.AddCommand(listCmd)
	projectCmd.AddCommand(variableCmd)
	projectCmd.AddCommand(variablesCmd)
}
