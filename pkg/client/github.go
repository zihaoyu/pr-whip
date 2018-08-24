package client

import (
	"context"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

// GenericGithubAPIClient is the github client interface
type GenericGithubAPIClient interface {
	// ListPullRequestsInRepo lists all the pull requests in a give repo
	ListPullRequestsInRepo(owner, repo string) ([]*github.PullRequest, error)
}

// GithubAPIClient is an implementation of the interface
type GithubAPIClient struct {
	ctx    context.Context
	client *github.Client
}

// NewGithubAPIClient creates a new github api client
func NewGithubAPIClient(ctx context.Context, client *github.Client) GenericGithubAPIClient {
	return &GithubAPIClient{
		ctx:    ctx,
		client: client,
	}
}

// ListPullRequestsInRepo lists all the pull requests in a give repo
func (c *GithubAPIClient) ListPullRequestsInRepo(owner, repo string) ([]*github.PullRequest, error) {
	opts := github.PullRequestListOptions{
		State: "open",
	}
	pulls, _, err := c.client.PullRequests.List(c.ctx, owner, repo, &opts)
	log.Infof("fetched %d pull requests from repo %s/%s", len(pulls), owner, repo)
	return pulls, err
}
