package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/jasonlvhit/gocron"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	client "github.com/zihaoyu/pr-whip/pkg/client"
	cfg "github.com/zihaoyu/pr-whip/pkg/config"
	wkr "github.com/zihaoyu/pr-whip/pkg/worker"

	"golang.org/x/oauth2"
)

var (
	githubAPIKey string
	slackAPIKey  string
	configFile   string
)

func main() {
	setLogging()

	flag.StringVar(&githubAPIKey, "github-api-key", "", "Github API key")
	flag.StringVar(&slackAPIKey, "slack-api-key", "", "Slack API key")
	flag.StringVar(&configFile, "config-file", "", "path to rules config file")
	flag.Parse()

	if githubAPIKey == "" {
		log.Fatal("github api key must be provided")
	}
	if slackAPIKey == "" {
		log.Fatal("slack api key must be provided")
	}
	if configFile == "" {
		log.Fatal("config file path must be provided")
	}

	config, err := cfg.FromFile(configFile)
	if err != nil {
		log.Fatal("config file cannot be loaded")
	}

	log.Infof("loaded %d rules", len(config.Rules))

	ctx := context.Background()

	githubHTTPClient := createGithubHTTPClient(ctx, githubAPIKey)
	slackHTTPClient := createSlackHTTPClient(slackAPIKey)

	githubAPIClient := client.NewGithubAPIClient(ctx, githubHTTPClient)
	slackAPIClient := client.NewSlackAPIClient(slackHTTPClient)
	createSchedules(config.Rules, githubAPIClient, slackAPIClient)
	// function Start start all the pending jobs
	<-gocron.Start()
}

func setLogging() {
	log.SetFormatter(&log.JSONFormatter{})
}

func createGithubHTTPClient(ctx context.Context, apiKey string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: apiKey,
		},
	)
	httpClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(httpClient)
	return client
}

func createSlackHTTPClient(apiKey string) *slack.Client {
	return slack.New(apiKey)
}

func createSchedules(rules []cfg.Rule, github client.GenericGithubAPIClient, slack client.GenericSlackAPIClient) {
	worker := wkr.New(github, slack)
	gocron.ChangeLoc(time.UTC)
	for _, rule := range rules {
		for _, schedule := range rule.Schedules {
			parts := strings.Split(schedule, " ")
			if len(parts) < 2 {
				log.Errorf("invalid schedule %s", schedule)
				continue
			}

			days := strings.Split(parts[0], ",")
			times := strings.Split(parts[1], ",")

			for _, t := range times {
				for _, d := range days {
					j := gocron.Every(1)
					j, err := applyDayOfWeek(j, d)
					if err != nil {
						continue
					}
					j.At(t).Do(worker.Do, rule)
				}
			}
		}
	}
}

func applyDayOfWeek(job *gocron.Job, day string) (*gocron.Job, error) {
	switch day {
	case "Monday":
		return job.Monday(), nil
	case "Tuesday":
		return job.Tuesday(), nil
	case "Wednesday":
		return job.Wednesday(), nil
	case "Thursday":
		return job.Thursday(), nil
	case "Friday":
		return job.Friday(), nil
	case "Saturday":
		return job.Saturday(), nil
	case "Sunday":
		return job.Sunday(), nil
	default:
		log.Errorf("invalid day specifier: %s", day)
		return job, fmt.Errorf("invalid day specifier: %s", day)
	}
}
