package github_global

import (
	"context"
	"errors"
	"fmt"
	"github.com/wtfutil/wtf/logger"
	"net/http"
	"net/url"
	"strings"

	ghb "github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

const (
	pullRequestsPath = "/pulls"
	issuesPath       = "/issues"

	myPullRequestsQuery    = "is:open is:pr author:%s archived:false"
	myReviewRequestedQuery = "is:open is:pr review-requested:%s archived:false"
	myAssignedIssuesQuery  = "is:open is:issue assignee:%s archived:false "
)

// data defines a new GitHub Repo structure
type data struct {
	apiKey    string
	username  string
	baseURL   string
	uploadURL string
	sortOrder string

	MyPullRequests   []*ghb.Issue
	MyReviewRequests []*ghb.Issue
	MyAssignedIssues []*ghb.Issue
	Err              error
}

// newDate returns a new Github Repo with a name, owner, apiKey, baseURL and uploadURL
func newData(settings *Settings) *data {
	repo := data{
		apiKey:    settings.apiKey,
		username:  settings.username,
		baseURL:   settings.baseURL,
		uploadURL: settings.uploadURL,
		sortOrder: settings.sortOrder,
	}

	return &repo
}

// OpenPulls will open the GitHub Pull Requests URL using the utils helper
func (data *data) OpenPulls() {
	panic(errors.New("implement open global PR search page"))
	//utils.OpenFile(*data.RemoteRepo.HTMLURL + pullRequestsPath)
}

// OpenIssues will open the GitHub Issues URL using the utils helper
func (data *data) OpenIssues() {
	panic(errors.New("implement open global issues search page"))
	//utils.OpenFile(*data.RemoteRepo.HTMLURL + issuesPath)
}

// Refresh reloads the github data via the Github API
func (data *data) Refresh() {
	data.MyPullRequests, data.Err = data.loadDefaultSearch(myPullRequestsQuery)
	if data.Err != nil {
		return
	}
	data.MyReviewRequests, data.Err = data.loadDefaultSearch(myReviewRequestedQuery)
	if data.Err != nil {
		return
	}
	data.MyAssignedIssues, data.Err = data.loadDefaultSearch(myAssignedIssuesQuery)
	if data.Err != nil {
		return
	}
}

/* -------------------- Counts -------------------- */

// IssueCount return the total amount of issues as an int
func (data *data) IssueCount() int {
	return len(data.MyAssignedIssues)
}

// PullRequestCount returns the total amount of pull requests as an int
func (data *data) PullRequestCount() int {
	return len(data.MyPullRequests)
}

/* -------------------- Unexported Functions -------------------- */

func (data *data) isGitHubEnterprise() bool {
	if len(data.baseURL) > 0 {
		if data.uploadURL == "" {
			data.uploadURL = data.baseURL
		}
		return true
	}
	return false
}

func (data *data) oauthClient() *http.Client {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: data.apiKey},
	)

	return oauth2.NewClient(context.Background(), tokenService)
}

func (data *data) githubClient() (*ghb.Client, error) {
	oauthClient := data.oauthClient()

	if data.isGitHubEnterprise() {
		return ghb.NewEnterpriseClient(data.baseURL, data.uploadURL, oauthClient)
	}

	return ghb.NewClient(oauthClient), nil
}

// myPullRequests returns a list of pull requests created by username on this repo
func (data *data) myPullRequests() []*ghb.PullRequest {
	return data.issuesToPrs(data.MyPullRequests)
}

// myReviewRequests returns a list of pull requests for which username has been
// requested to do a code review
func (data *data) myReviewRequests() []*ghb.PullRequest {
	return data.issuesToPrs(data.MyReviewRequests)
}

func (data *data) issuesToPrs(issues []*ghb.Issue) []*ghb.PullRequest {
	github, err := data.githubClient()
	if err != nil {
		return nil
	}

	ret := make([]*ghb.PullRequest, 0, len(issues))
	for _, issue := range issues {
		if !issue.IsPullRequest() {
			continue
		}
		owner, repo := parseRepositoryUrl(issue.GetRepositoryURL())
		pr, _, err := github.PullRequests.Get(context.Background(), owner, repo, issue.GetNumber())
		if err != nil || pr == nil {
			logger.LogJson("Could not resolve PR", issue)
		} else {
			ret = append(ret, pr)
		}
	}
	return ret
}

func (data *data) customIssueQuery(filter string, perPage int) *ghb.IssuesSearchResult {
	github, err := data.githubClient()
	if err != nil {
		return nil
	}

	opts := &ghb.SearchOptions{}
	if perPage != 0 {
		opts.ListOptions.PerPage = perPage
	}

	prs, _, _ := github.Search.Issues(context.Background(), filter, opts)
	return prs
}

func (data *data) loadDefaultSearch(search string) ([]*ghb.Issue, error) {
	github, err := data.githubClient()
	if err != nil {
		return nil, err
	}

	opts := new(ghb.SearchOptions)
	opts.Sort = data.sortOrder
	opts.ListOptions.PerPage = 5

	prIssues, _, err := github.Search.Issues(context.Background(), fmt.Sprintf(search, data.username), opts)

	if err != nil {
		return nil, err
	}

	return prIssues.Issues, nil
}

func parseRepositoryUrl(repoUrl string) (owner, name string) {
	url, _ := url.Parse(repoUrl)
	path := strings.Split(url.Path, "/")

	return path[len(path)-2], path[len(path)-1]
}
