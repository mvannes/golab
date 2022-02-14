package mergerequest

import (
	"log"
	"strings"

	"github.com/mvannes/golab/gitlab"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	lab "github.com/xanzy/go-gitlab"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "List all open merge requests",
	Run: func(cmd *cobra.Command, args []string) {
		c := gitlab.NewClient()

		mrs, err := c.MergeRequests(gitlab.Open)
		if nil != err {
			log.Fatal(err)
		}

		headerFmt := color.New(color.BgBlue, color.Underline).SprintfFunc()
		searchTerms := viper.GetStringSlice("merge-request-search-words")
		tbl := table.New("Title", "Author", "Url", "Updated")
		tbl.WithHeaderFormatter(headerFmt)
		for _, mr := range mrs {
			if !matchesSearchTerms(mr, searchTerms) {
				continue
			}
			tbl.AddRow(mr.Title, mr.Author.Username, mr.WebURL, mr.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
		tbl.Print()
	},
}

func checkUserForTerm(u lab.BasicUser, term string) bool {
	if strings.Contains(u.Name, term) {
		return true
	}
	if strings.Contains(u.Username, term) {
		return true
	}
	return false
}

func matchesSearchTerms(mr lab.MergeRequest, searchTerms []string) bool {
	var users []*lab.BasicUser
	users = append(users, mr.Author)
	users = append(users, mr.Assignees...)

	for _, t := range searchTerms {
		for _, u := range users {
			if checkUserForTerm(*u, t) {
				return true
			}
		}

		if strings.Contains(mr.Description, t) {
			return true
		}
	}

	return false
}

func init() {
	mergeRequestCmd.AddCommand(openCmd)
}
