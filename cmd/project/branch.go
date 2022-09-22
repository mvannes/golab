package project

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/mvannes/golab/gitlab"
	"github.com/spf13/cobra"
)

var getBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Get a branch of the given project [namespace] [project-name] [branch-name]",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		project, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		if nil == project {
			log.Fatal(errors.New("No project found"))
		}

		branch, err := c.Branch(*project, args[2])
		if nil != err {
			log.Fatal(err.Error())
		}
		fmt.Println(branch.Name)
		fmt.Println(branch.Protected)
	},
}

var getBranchListCmd = &cobra.Command{
	Use:   "branches",
	Short: "Get a branch of the given project [namespace] [project-name]",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		project, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		if nil == project {
			log.Fatal(errors.New("No project found"))
		}

		branches, err := c.Branches(*project)
		if nil != err {
			log.Fatal(err.Error())
		}

		for _, b := range branches {
			fmt.Println(b.Name)
		}
	},
}

var protectCmd = &cobra.Command{
	Use:   "protect-branch",
	Short: "Protect a branch of the given project",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		project, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		if nil == project {
			log.Fatal(errors.New("No project found"))
		}

		branch, err := c.Branch(*project, args[2])
		if nil != err {
			log.Fatal(err.Error())
		}

		c.ProtectBranch(*project, *branch)
	},
}

var unprotectCmd = &cobra.Command{
	Use:   "unprotect-branch",
	Short: "Unprotect a branch of the given project",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		project, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		if nil == project {
			log.Fatal(errors.New("No project found"))
		}

		branch, err := c.Branch(*project, args[2])
		if nil != err {
			log.Fatal(err.Error())
		}

		c.UnprotectBranch(*project, *branch)
	},
}

var pruneStaleBranchesCmd = &cobra.Command{
	Use:   "prune-stale-branches",
	Short: "Prune stale branches of the given project [namespace] [project-name]",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		project, err := c.Project(args[0], args[1])
		if nil != err {
			log.Fatal(err.Error())
		}
		if nil == project {
			log.Fatal(errors.New("No project found"))
		}

		branches, err := c.Branches(*project)
		if nil != err {
			log.Fatal(err.Error())
		}

		for _, b := range branches {
			if b.Default == true {
				continue
			}
			cmt := b.Commit

			p := promptui.Select{
				Label: fmt.Sprint("Remove branch ", b.Name, " last commited to at", cmt.CommittedDate.String(), " "),
				Items: []string{"yes", "no"},
			}

			_, i, err := p.Run()
			if nil != err {
				log.Fatal(err)
			}
			if i == "no" {
				continue
			}

			if err = c.RemoveBranch(*project, b); nil != err {
				log.Fatal(err)
			}
			fmt.Println("Branch ", b.Name, " removed")
		}
	},
}

var unprotectedDefaultBranchesCmd = &cobra.Command{
	Use:   "unprotected-default-branches",
	Short: "List any default branches for projects in [namespace]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()
		projects, err := c.Projects(args[0])
		if nil != err {
			log.Fatal(err.Error())
		}

		for _, project := range projects {
			if project.Archived == true {
				continue
			}
			branches, err := c.Branches(project)
			if nil != err {
				log.Fatal(err.Error())
			}

			for _, b := range branches {
				if b.Default == true && b.Protected == false {
					fmt.Println(project.PathWithNamespace)
				}
			}
		}
	},
}

func init() {
	projectCmd.AddCommand(getBranchListCmd)
	projectCmd.AddCommand(getBranchCmd)
	projectCmd.AddCommand(protectCmd)
	projectCmd.AddCommand(unprotectCmd)
	projectCmd.AddCommand(pruneStaleBranchesCmd)
	projectCmd.AddCommand(unprotectedDefaultBranchesCmd)
}
