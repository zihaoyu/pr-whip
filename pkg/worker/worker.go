package worker

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"

	"github.com/zihaoyu/pr-whip/pkg/client"
	cfg "github.com/zihaoyu/pr-whip/pkg/config"
)

const (
	slackMessageUsername = "PR Whip"
	slackMessageIconURL  = "http://getdrawings.com/image/whip-drawing-60.png"
	slackMessageText     = "Some pull requests are open. Please review."
)

// Worker does the work according to given rule
type Worker struct {
	github client.GenericGithubAPIClient
	slack  client.GenericSlackAPIClient
}

// New creates a new worker with service clients
func New(githubClient client.GenericGithubAPIClient, slackClient client.GenericSlackAPIClient) *Worker {
	return &Worker{
		github: githubClient,
		slack:  slackClient,
	}
}

// Do runs the actual work
func (w *Worker) Do(rule cfg.Rule) {
	pulls := w.listPullRequests(rule)
	if len(pulls) > 0 {
		msg := buildMessage(pulls)
		for _, channel := range rule.Channels {
			err := w.slack.Notify(channel, slackMessageText, msg)
			if err != nil {
				log.Errorf("error sending message to slack channel %s: %v", channel, err)
			}
		}
	} else {
		log.Info("no pull request to report")
	}

}

func buildMessage(pulls []*github.PullRequest) slack.PostMessageParameters {
	params := slack.NewPostMessageParameters()
	var attachments []slack.Attachment

	// create attachment for each pull request
	for _, pull := range pulls {
		age := time.Now().UTC().Sub(*pull.CreatedAt)
		age = age.Round(time.Hour)
		var text string
		if age.Hours() < 24.0 {
			text = fmt.Sprintf("%d hours", int64(age.Hours()))
		} else {
			h := math.Mod(age.Hours(), 24)
			text = fmt.Sprintf("%d days %d hours", int64(age.Hours()/24), int64(h))
		}

		field := slack.AttachmentField{
			Title: "Age",
			Value: text,
			Short: true,
		}
		attachment := slack.Attachment{
			Fallback:  *pull.URL,
			Color:     "#36a64f",
			Title:     fmt.Sprintf("%s#%d", *pull.Base.Repo.FullName, *pull.Number),
			TitleLink: *pull.URL,
			Text:      *pull.Title,
			Fields:    []slack.AttachmentField{field},
		}

		attachments = append(attachments, attachment)
	}
	params.Attachments = attachments
	params.Username = slackMessageUsername
	params.IconURL = slackMessageIconURL

	return params
}

func (w *Worker) listPullRequests(rule cfg.Rule) []*github.PullRequest {
	var pulls []*github.PullRequest
	for _, repo := range rule.Repos {
		parts := strings.SplitN(repo, "/", 2)
		if len(parts) < 2 {
			log.Errorf("invalid repo identifier: %s", repo)
			continue
		}

		ps, err := w.github.ListPullRequestsInRepo(parts[0], parts[1])
		if err != nil {
			log.Errorf("error fetching pull requests from repo %s: %v", repo, err)
			continue
		}

		for _, p := range ps {
			log.Infof("checking pull request %s#%d", repo, *p.Number)

			pullRequestAge := time.Now().UTC().Sub(*p.CreatedAt)

			ignored := false

			// check if pull request is too young
			youngerThan := rule.Ignore.YoungerThan
			if youngerThan != "" {
				youngThreshold, err := time.ParseDuration(youngerThan)
				if err != nil {
					log.Errorf("error parsing ignore lower bound %s: %v", youngerThan, err)
				} else {
					if pullRequestAge < youngThreshold {
						log.Infof("pull request age %s is younger than threshold %s, ignore", pullRequestAge.String(), youngerThan)
						ignored = true
					}
				}
			}

			// check if pull request is too old if i is not too young
			olderThan := rule.Ignore.OlderThan
			if !ignored && olderThan != "" {
				oldThreshold, err := time.ParseDuration(olderThan)
				if err != nil {
					log.Errorf("error parsing ignore upper bound %s: %v", olderThan, err)
				} else {
					if pullRequestAge > oldThreshold {
						log.Infof("pull request age %s is older than threshold %s, ignore", pullRequestAge.String(), olderThan)
						ignored = true
					}
				}
			}

			// if all good, add pull request to list
			if !ignored {
				pulls = append(pulls, p)
			}
		}
	}
	return pulls
}
