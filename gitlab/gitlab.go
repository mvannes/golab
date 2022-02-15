package gitlab

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	gitlab *gitlab.Client
}
type MergeRequestState string

const (
	Open   MergeRequestState = "opened"
	Closed MergeRequestState = "closed"
	Locked MergeRequestState = "locked"
	Merged MergeRequestState = "merged"
	All    MergeRequestState = "all"
)

var client *GitlabClient
var projWg sync.WaitGroup
var branchWg sync.WaitGroup
var mrWg sync.WaitGroup

func NewClient() *GitlabClient {
	if client != nil {
		return client
	}
	gitlabToken := viper.GetString("gitlab-token")
	if "" == gitlabToken {
		log.Fatal(errors.New("No gitlab token configured"))
	}
	gitlabBaseUrl := viper.GetString("gitlab-base-url")

	if "" == gitlabBaseUrl {
		log.Fatal(errors.New("No gitlab base url configured"))
	}
	gitlabClient, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabBaseUrl))

	if nil != err {
		log.Fatal(err)
	}
	client = &GitlabClient{
		gitlab: gitlabClient,
	}

	return client
}

func (g *GitlabClient) Projects(namespace string) ([]gitlab.Project, error) {
	c := make(chan gitlab.Project)

	_, r, err := g.gitlab.Projects.ListProjects(&gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{PerPage: 50, Page: 1}})
	if nil != err {
		return make([]gitlab.Project, 0), err
	}

	for i := 1; i <= r.TotalPages; i++ {
		projWg.Add(1)
		go doProjectListRequest(*g.gitlab, c, i)
	}
	go func() {
		projWg.Wait()
		close(c)
	}()

	var result []gitlab.Project
	for p := range c {
		if p.Namespace.Name == namespace {
			result = append(result, p)
		}
	}

	return result, nil
}

func doProjectListRequest(c gitlab.Client, projectChan chan gitlab.Project, page int) {
	defer projWg.Done()
	projects, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{PerPage: 50, Page: page}})

	if nil != err {
		log.Fatal(err)
	}
	for _, p := range projects {
		projectChan <- *p
	}
}

func (g *GitlabClient) Project(namespace, name string) (*gitlab.Project, error) {
	p, _, err := g.gitlab.Projects.GetProject(fmt.Sprint(namespace, "/", name), &gitlab.GetProjectOptions{})
	return p, err
}

func (g *GitlabClient) Branches(p gitlab.Project) ([]gitlab.Branch, error) {
	_, r, err := g.gitlab.Branches.ListBranches(p.ID, &gitlab.ListBranchesOptions{ListOptions: gitlab.ListOptions{PerPage: 50, Page: 1}})
	if nil != err {
		return make([]gitlab.Branch, 0), err
	}
	c := make(chan gitlab.Branch)

	for i := 1; i <= r.TotalPages; i++ {
		branchWg.Add(1)
		go doBranchListRequest(*g.gitlab, c, p, i)
	}
	go func() {
		branchWg.Wait()
		close(c)
	}()

	var result []gitlab.Branch
	for b := range c {
		result = append(result, b)
	}
	return result, nil
}

func doBranchListRequest(c gitlab.Client, branchChan chan gitlab.Branch, p gitlab.Project, page int) {
	defer branchWg.Done()
	branches, _, err := c.Branches.ListBranches(
		p.ID,
		&gitlab.ListBranchesOptions{ListOptions: gitlab.ListOptions{PerPage: 50, Page: page}},
	)
	if nil != err {
		log.Fatal(err)
	}
	for _, b := range branches {
		branchChan <- *b
	}
}

func (g *GitlabClient) Branch(p gitlab.Project, branchName string) (*gitlab.Branch, error) {
	branch, _, err := g.gitlab.Branches.GetBranch(p.ID, branchName)
	return branch, err
}

func (g *GitlabClient) ProtectBranch(p gitlab.Project, b gitlab.Branch) error {
	_, _, err := g.gitlab.Branches.ProtectBranch(p.ID, b.Name, &gitlab.ProtectBranchOptions{})
	return err
}

func (g *GitlabClient) UnprotectBranch(p gitlab.Project, b gitlab.Branch) error {
	_, _, err := g.gitlab.Branches.UnprotectBranch(p.ID, b.Name)
	return err
}

func (g *GitlabClient) MergeRequests(state MergeRequestState) ([]gitlab.MergeRequest, error) {
	scopeOpt := "all"
	opts := gitlab.ListMergeRequestsOptions{Scope: &scopeOpt, ListOptions: gitlab.ListOptions{PerPage: 50, Page: 1}}
	if state != All {
		optState := string(state)
		opts.State = &optState
	}

	_, r, err := g.gitlab.MergeRequests.ListMergeRequests(&opts)
	if err != nil {
		return make([]gitlab.MergeRequest, 0), err
	}

	c := make(chan gitlab.MergeRequest)

	for i := 1; i <= r.TotalPages; i++ {
		mrWg.Add(1)
		go doMergeRequestListRequest(*g.gitlab, c, opts, i)
	}
	go func() {
		mrWg.Wait()
		close(c)
	}()

	var result []gitlab.MergeRequest
	for mr := range c {
		result = append(result, mr)
	}
	return result, nil

}

func doMergeRequestListRequest(c gitlab.Client, mrChan chan gitlab.MergeRequest, opts gitlab.ListMergeRequestsOptions, page int) {
	defer mrWg.Done()

	opts.Page = page
	mrs, _, err := c.MergeRequests.ListMergeRequests(&opts)
	if nil != err {
		log.Fatal(err)
	}
	for _, mr := range mrs {
		mrChan <- *mr
	}
}
